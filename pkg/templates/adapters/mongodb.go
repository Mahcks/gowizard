package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
)

type MongoDBAdapter struct {
	name             string // name of the adapter
	*domain.Settings        // settings of the project
}

// GetName returns the name of the adapter
func (a *MongoDBAdapter) GetName() string {
	return a.name
}

func NewMongoDBAdapter(settings *domain.Settings) domain.ModuleI {
	return &MongoDBAdapter{
		name:     "mongodb",
		Settings: settings,
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (m *MongoDBAdapter) ConfigGo() *j.Statement {
	return j.Id("MongoDB").Struct(
		j.Id("URI").String().Tag(map[string]string{"mapstructure": "uri", "json": "uri"}),
	).Tag(map[string]string{"mapstructure": "mongodb", "json": "mongodb"})
}

// ConfigGo is the configuration of the adapter in Go format
func (m *MongoDBAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"mongodb": map[string]interface{}{
			"uri": "mongodb://localhost:27017",
		},
	}
}

// AppInit is the code that will be added to the END internal/app/app.go Run() function
func (m *MongoDBAdapter) AppInit() []j.Code {
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

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (m *MongoDBAdapter) AppShutdown() []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("mongodb").Dot("Close").Call(j.Id("gCtx")),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (a *MongoDBAdapter) Service() *j.File {
	f := j.NewFilePathName(a.Settings.Module+"/pkg/mongodb", "mongodb")

	// Service struct
	ptr := j.Op("*")
	sStruct := j.Type().Id("MongoDB").Struct(
		j.Id("ctx").Qual("context", "Context"),
		j.Id("Client").Add(ptr).Qual("go.mongodb.org/mongo-driver/mongo", "Client"),
	)

	f.Add(sStruct)

	f.Func().Id("New").Params(
		j.Id("gCtx").Qual("context", "Context"),
		j.Id("uri").String(),
	).Op("(").List(j.Op("*").Add(j.Id("MongoDB"), j.Op(","), j.Error()).Op(")")).Block(
		j.List(j.Id("client"), j.Err()).Op(":=").Qual("go.mongodb.org/mongo-driver/mongo", "Connect").Params(
			j.Id("gCtx"),
			j.Qual("go.mongodb.org/mongo-driver/mongo/options", "Client").Call().Dot("ApplyURI").Call(j.Id("uri")),
		),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Err()),
		),
		j.Line(),
		j.Line().Comment("Ping to see if connection was successful"),
		j.Err().Op("=").Id("client").Dot("Ping").Call(j.Id("gCtx"), j.Qual("go.mongodb.org/mongo-driver/mongo/readpref", "Primary").Call()),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Err()),
		),
		j.Line(),
		j.Return(j.Op("&").Id("MongoDB").Values(j.Dict{
			j.Id("Client"): j.Id("client"),
		}), j.Nil()),
	)

	f.Add(j.Line())

	// Close method
	f.Func().Params(j.Id("m").Op("*").Id("MongoDB")).Id("Close").Params(
		j.Id("gCtx").Qual("context", "Context"),
	).Block(
		j.If(j.Id("m.Client").Op("!=").Nil()).Block(
			j.Id("m.Client").Dot("Disconnect").Call(j.Id("gCtx")),
		),
	)

	return f
}
