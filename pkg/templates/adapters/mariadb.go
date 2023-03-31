package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
)

type MariaDBAdapter struct {
	name string
	*domain.Settings
}

func (adp *MariaDBAdapter) GetName() string {
	return adp.name
}

func NewMariaDBAdapter(settings *domain.Settings) domain.ModuleI {
	return &MariaDBAdapter{
		name:     "mariadb",
		Settings: settings,
	}
}

func (adp *MariaDBAdapter) ConfigGo() *j.Statement {
	return j.Id("MariaDB").Struct(
		j.Id("Host").String().Tag(map[string]string{"mapstructure": "host", "json": "host"}),
		j.Id("Port").String().Tag(map[string]string{"mapstructure": "port", "json": "port"}),
		j.Id("Username").String().Tag(map[string]string{"mapstructure": "username", "json": "username"}),
		j.Id("Password").String().Tag(map[string]string{"mapstructure": "password", "json": "password"}),
		j.Id("Database").String().Tag(map[string]string{"mapstructure": "database", "json": "database"}),
	).Tag(map[string]string{"mapstructure": "mariadb", "json": "mariadb"})
}

func (m *MariaDBAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"mariadb": map[string]interface{}{
			"host":     "localhost",
			"port":     "3306",
			"username": "user",
			"password": "password",
			"database": "testdb",
		},
	}
}

func (adp *MariaDBAdapter) AppInit() []j.Code {
	return []j.Code{
		j.List(j.Id("mdb"), j.Err()).Op(":=").Qual(adp.Module+"/pkg/mariadb", "New").Params(j.Id("cfg.MariaDB.Host"), j.Id("cfg.MariaDB.Port"), j.Id("cfg.MariaDB.Database"), j.Id("cfg.MariaDB.Username"), j.Id("cfg.MariaDB.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - mariadb.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("MariaDB initialized")),
		j.Line(),
	}
}

func (adp *MariaDBAdapter) AppShutdown() []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("mdb").Dot("Close").Call(),
	}
}

func (adp *MariaDBAdapter) Service() *j.File {
	f := j.NewFilePathName(adp.Settings.Module+"/pkg/mariadb", "mariadb")

	// Service struct
	ptr := j.Op("*")
	sStruct := j.Type().Id("MariaDB").Struct(
		j.Id("DB").Add(ptr).Qual("database/sql", "DB"),
	)

	f.Add(sStruct)

	// New function
	f.Func().Id("New").Params(
		j.Id("host"),
		j.Id("port"),
		j.Id("database"),
		j.Id("username"),
		j.Id("password").String(),
	).Op("(").List(j.Op("*").Add(j.Id("MariaDB"), j.Op(","), j.Error()).Op(")")).Block(
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
		j.Return(j.Op("&").Id("MariaDB").Values(j.Dict{
			j.Id("DB"): j.Id("client"),
		}), j.Nil()),
	)

	f.Add(j.Line())

	// Close function
	f.Func().Params(j.Id("m").Op("*").Id("MariaDB")).Id("Close").Params().Error().Block(
		j.If(
			j.Id("m").Dot("DB").Op("!=").Nil().Block(
				j.Id("m").Dot("DB").Dot("Close").Call(),
			),
		),
		j.Line(),
		j.Return(j.Nil()),
	)

	return f
}
