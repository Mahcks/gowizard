package mariadb

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/internal/domain"
	"gopkg.in/yaml.v2"
)

type Adapter struct {
	name             string // name of the adapter
	*domain.Settings        // settings of the project
}

// GetName returns the name of the adapter
func (a *Adapter) GetName() string {
	return a.name
}

func NewAdapter(name string, settings *domain.Settings) domain.ModuleI {
	return &Adapter{
		name:     name,
		Settings: settings,
	}
}

func (m *Adapter) ConfigGo() *j.Statement {
	return j.Id("MariaDB").Struct(
		j.Id("Host").String().Tag(map[string]string{"mapstructure": "host", "json": "host"}),
		j.Id("Port").String().Tag(map[string]string{"mapstructure": "port", "json": "port"}),
		j.Id("Username").String().Tag(map[string]string{"mapstructure": "username", "json": "username"}),
		j.Id("Password").String().Tag(map[string]string{"mapstructure": "password", "json": "password"}),
		j.Id("Database").String().Tag(map[string]string{"mapstructure": "database", "json": "database"}),
	).Tag(map[string]string{"mapstructure": "mariadb", "json": "mariadb"})
}

func (m *Adapter) ConfigYAML() ([]byte, error) {
	data := map[string]interface{}{
		"mariadb": map[string]interface{}{
			"host":     "localhost",
			"port":     "3306",
			"username": "user",
			"password": "password",
			"database": "testdb",
		},
	}

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

func (m *Adapter) AppInit() []j.Code {
	return []j.Code{
		j.List(j.Id("mdb"), j.Err()).Op(":=").Qual(m.Module+"/pkg/mariadb", "New").Params(j.Id("cfg.MariaDB.Host"), j.Id("cfg.MariaDB.Port"), j.Id("cfg.MariaDB.Database"), j.Id("cfg.MariaDB.Username"), j.Id("cfg.MariaDB.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - mariadb.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("MariaDB initialized")),
	}
}

func (m *Adapter) AppShutdown() []j.Code {
	return []j.Code{
		j.Id("mdb").Dot("Close").Call(),
	}
}
