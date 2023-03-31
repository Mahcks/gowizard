package mariadb

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/internal/domain"
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
	return j.Id("MongoDB").Struct(
		j.Id("URI").String().Tag(map[string]string{"mapstructure": "uri", "json": "uri"}),
	).Tag(map[string]string{"mapstructure": "mongodb", "json": "mongodb"})
}

func (m *Adapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"mongodb": map[string]interface{}{
			"uri": "mongodb://localhost:27017",
		},
	}
}

func (m *Adapter) AppInit() []j.Code {
	return []j.Code{
		j.Line(),
		j.List(j.Id("mongodb"), j.Err()).Op(":=").Qual(m.Module+"/pkg/mongodb", "New").Params(j.Id("gCtx"), j.Id("cfg.MongoDB.URI")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - mongodb.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("MongoDB initialized")),
		j.Line(),
	}
}

func (m *Adapter) AppShutdown() []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("mongodb").Dot("Close").Call(j.Id("gCtx")),
	}
}
