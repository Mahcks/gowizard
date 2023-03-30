package mariadb

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/internal/domain"
)

type Adapter struct{}

func NewAdapter() domain.ModuleI {
	return &Adapter{}
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

func (m *Adapter) AppInit() []j.Code {
	return []j.Code{
		j.List(j.Id("mdb"), j.Err()).Op(":=").Qual("github.com/mahcks/test-project/pkg/mariadb", "New").Params(j.Id("cfg.MariaDB.Host"), j.Id("cfg.MariaDB.Port"), j.Id("cfg.MariaDB.Database"), j.Id("cfg.MariaDB.Username"), j.Id("cfg.MariaDB.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - mariadb.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("MariaDB initialized")),
	}
}

func (m *Adapter) AppShutdown() *j.Statement {
	return j.Id("mdb").Dot("Close").Call()
}
