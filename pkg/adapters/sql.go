package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type SQLAdapter struct {
	name        string // name of the adapter
	displayName string // name of the adapter that will be displayed in the CLI
}

// GetName returns the name of the adapter
func (adp *SQLAdapter) GetName() string {
	return adp.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (adp *SQLAdapter) GetDisplayName() string {
	return adp.displayName
}

func NewSQLAdapter() domain.ModuleI {
	return &SQLAdapter{
		name:        "sql",
		displayName: "SQL",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (m *SQLAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"sql": map[string]interface{}{
			"host":     "localhost",
			"port":     "3306",
			"username": "user",
			"password": "password",
			"database": "testdb",
		},
	}
}

// ConfigGo is the configuration of the adapter in Go format
func (adp *SQLAdapter) ConfigGo() *j.Statement {
	return j.Id("SQL").Struct(
		j.Id("Host").String().Tag(map[string]string{"mapstructure": "host", "json": "host"}),
		j.Id("Port").String().Tag(map[string]string{"mapstructure": "port", "json": "port"}),
		j.Id("Username").String().Tag(map[string]string{"mapstructure": "username", "json": "username"}),
		j.Id("Password").String().Tag(map[string]string{"mapstructure": "password", "json": "password"}),
		j.Id("Database").String().Tag(map[string]string{"mapstructure": "database", "json": "database"}),
	).Tag(map[string]string{"mapstructure": "sql", "json": "sql"})
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (adp *SQLAdapter) AppInit(module string) []j.Code {
	return []j.Code{
		j.List(j.Id("db"), j.Err()).Op(":=").Qual(module+"/pkg/sql", "New").Params(j.Id("cfg.SQL.Host"), j.Id("cfg.SQL.Port"), j.Id("cfg.SQL.Database"), j.Id("cfg.SQL.Username"), j.Id("cfg.SQL.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("fmt", "Println").Call(j.Lit("error connecting to sql"), j.Err()),
		),
		j.Line(),
		j.Line(),
		j.Qual("fmt", "Println").Call(j.Lit("connected to sql")),
		j.Line(),
	}
}

// AppSelect - Each AppSelect branch is apart of a bigger switch statement that's in the internal/app/app.go Run() function
func (adp *SQLAdapter) AppSelect(module string) j.Code {

	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (adp *SQLAdapter) AppShutdown(module string) []j.Code {
	return []j.Code{
		j.Id("db").Dot("Close").Call(),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (adp *SQLAdapter) Service(module, path string) *j.File {
	f := j.NewFilePathName(module+"/pkg/sql", "sql")

	// Service struct
	sStruct := j.Type().Id("SQL").Struct(
		j.Id("DB").Add(utils.Jptr).Qual("database/sql", "DB"),
	)

	f.Add(sStruct)

	// New function
	f.Func().Id("New").Params(
		j.Id("host"),
		j.Id("port"),
		j.Id("database"),
		j.Id("username"),
		j.Id("password").String(),
	).Op("(").List(j.Op("*").Add(j.Id("SQL"), j.Op(","), j.Error()).Op(")")).Block(
		j.Id("connectionString").Op(":=").Qual("fmt", "Sprintf").Params(
			j.Lit("%s:%s@tcp(%s:%s)/%s?parseTime=true"),
			j.Id("username"),
			j.Id("password"), j.Id("host"), j.Id("port"), j.Id("database"),
		),
		j.List(j.Id("client"), j.Id("err")).Op(":=").Qual("database/sql", "Open").Params(j.Lit("mysql"), j.Id("connectionString")),
		j.If(j.Id("err").Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Id("err")),
		),
		j.Line(),
		j.Line().Comment("Ping the database to check if the connection is alive"),
		j.If(j.Id("err").Op(":=").Id("client").Dot("Ping").Call(),
			j.Id("err").Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Id("err")),
		),
		j.Line(),
		j.Return(j.Op("&").Id("SQL").Values(j.Dict{
			j.Id("DB"): j.Id("client"),
		}), j.Nil()),
	)

	f.Add(j.Line())

	// Close function
	f.Func().Params(j.Id("m").Op("*").Id("SQL")).Id("Close").Params().Error().Block(
		j.If(
			j.Id("m").Dot("DB").Op("!=").Nil().Block(
				j.Id("m").Dot("DB").Dot("Close").Call(),
			),
		),
		j.Line(),
		j.Return(j.Nil()),
	)

	err := f.Save(path + "/pkg/" + adp.name + "/adapter.go")
	if err != nil {
		utils.PrintError("error saving file: %s", err)
		return nil
	}

	return f
}
