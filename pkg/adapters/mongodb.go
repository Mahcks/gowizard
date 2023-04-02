package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type MongoDBAdapter struct {
	name        string // name of the adapter
	displayName string // name of the adapter that will be displayed in the CLI
}

// GetName returns the name of the adapter
func (adp *MongoDBAdapter) GetName() string {
	return adp.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (adp *MongoDBAdapter) GetDisplayName() string {
	return adp.displayName
}

func NewMongoDBAdapter() domain.ModuleI {
	return &MongoDBAdapter{
		name:        "mongodb",
		displayName: "MongoDB",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (adp *MongoDBAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"mongodb": map[string]interface{}{
			"uri": "mongodb://localhost:27017",
		},
	}
}

// ConfigGo is the configuration of the adapter in Go format
func (adp *MongoDBAdapter) ConfigGo() *j.Statement {
	return j.Id("MongoDB").Struct(
		j.Id("URI").String().Tag(map[string]string{"mapstructure": "uri", "json": "uri"}),
	).Tag(map[string]string{"mapstructure": "mongodb", "json": "mongodb"})
}

// AppInit is the code that will be added to the END internal/app/app.go Run() function
func (adp *MongoDBAdapter) AppInit(module string) []j.Code {
	return []j.Code{
		j.Line(),
		j.List(j.Id("mongodb"), j.Err()).Op(":=").Qual(module+"/pkg/mongodb", "New").Params(j.Id("gCtx"), j.Id("cfg.MongoDB.URI")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - mongodb.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("MongoDB initialized")),
		j.Line(),
	}
}

// AppSelect - Each AppSelect branch is apart of a bigger switch statement that's in the internal/app/app.go Run() function
func (adp *MongoDBAdapter) AppSelect(module string) j.Code {

	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (adp *MongoDBAdapter) AppShutdown(module string) []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("mongodb").Dot("Close").Call(j.Id("gCtx")),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (adp *MongoDBAdapter) Service(module string) *j.File {
	f := j.NewFilePathName(module+"/pkg/mongodb", "mongodb")

	// Service struct
	sStruct := j.Type().Id("MongoDB").Struct(
		j.Id("ctx").Qual("context", "Context"),
		j.Id("Client").Add(utils.Jptr).Qual("go.mongodb.org/mongo-driver/mongo", "Client"),
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
