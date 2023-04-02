package services

import (
	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type RESTGin struct {
	name        string // name of the service
	displayName string // name of the adapter that will be displayed in the CLI
}

// GetName returns the name of the service
func (svc *RESTGin) GetName() string {
	return svc.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (svc *RESTGin) GetDisplayName() string {
	return "Gin"
}

func NewRESTGinService() domain.ModuleI {
	return &RESTGin{
		name:        "rest_gin",
		displayName: "Gin",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (m *RESTGin) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the adapter in Go format
func (svc *RESTGin) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (svc *RESTGin) AppInit(module string) []j.Code {
	return []j.Code{
		j.Id("handler").Op(":=").Qual("github.com/gin-gonic/gin", "New").Call(),
		j.Line(),
		j.Id("httpServer").Op(":=").Qual(module+"/pkg/rest_gin", "New").Call(j.Id("handler")),
	}
}

func (svc *RESTGin) AppSelect(module string) j.Code {
	return j.Case(
		j.Id("err").Op("=").Op("<-").Id("httpServer").Dot("Notify").Call()).Block(
		j.Qual("go.uber.org/zap", "S").Call().Dot("Errorw").Params(j.Lit("app - httpServer"), j.Lit("Notify()"), j.Id("err")),
	)
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (svc *RESTGin) AppShutdown(module string) []j.Code {
	return []j.Code{
		j.Id("err").Op("=").Id("httpServer").Dot("Shutdown").Call(),
		j.Line(),
		j.If(j.Id("err").Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Errorw").Params(j.Lit("app - httpServer"), j.Lit("Shutdown()"), j.Id("err")),
		),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (svc *RESTGin) Service(module string) *j.File {
	f := j.NewFilePathName(module+"/pkg/rest_gin", "httpserver")

	// Service struct
	sStruct := j.Type().Id("Service").Struct(
		j.Id("server").Add(utils.Jptr).Qual("net/http", "Server"),
		j.Id("notify").Chan().Error(),
		j.Id("shutdownTimeout").Qual("time", "Duration"),
	)

	f.Add(sStruct)

	f.Var().Id("defaultReadTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)
	f.Var().Id("defaultWriteTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)
	f.Var().Id("defaultAddr").Op("=").Lit(":80")
	f.Var().Id("defaultShutdownTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)

	// New service
	f.Func().Id("New").Params(j.Id("handler").Qual("net/http", "Handler")).Add(utils.Jptr).Id("Service").Block(
		j.Id("httpServer").Op(":=").Add(utils.Rptr).Qual("net/http", "Server").Values(j.Dict{
			j.Id("Handler"):      j.Id("handler"),
			j.Id("ReadTimeout"):  j.Id("defaultReadTimeout"),
			j.Id("WriteTimeout"): j.Id("defaultWriteTimeout"),
			j.Id("Addr"):         j.Id("defaultAddr"),
		}),
		j.Line(),
		j.Id("s").Op(":=").Add(utils.Rptr).Id("Service").Values(
			j.Dict{
				j.Id("server"):          j.Id("httpServer"),
				j.Id("notify"):          j.Make(j.Chan().Error().Op(",").Lit(1)),
				j.Id("shutdownTimeout"): j.Id("defaultShutdownTimeout"),
			},
		),
		j.Line(),
		j.Id("s").Dot("start").Call(),
		j.Line(),
		j.Return(j.Id("s")),
	)

	f.Line()

	// start()
	f.Func().Params(j.Id("s").Add(utils.Jptr).Id("Service")).Id("start").Params().Block(
		j.Id("go").Func().Params().Block(
			j.Id("s").Dot("notify").Op("<-").Id("s").Dot("server").Dot("ListenAndServe").Call(),
			j.Id("close").Call(j.Id("s").Dot("notify")),
		).Call(),
	)

	f.Line()

	// Notify()
	f.Func().Params(j.Id("s").Add(utils.Jptr).Id("Service")).Id("Notify").Params().Op("<-").Chan().Error().Block(
		j.Return(j.Id("s").Dot("notify")),
	)

	f.Line()

	// Shutdown()
	f.Func().Params(j.Id("s").Add(utils.Jptr).Id("Service")).Id("Shutdown").Params().Error().Block(
		j.List(j.Id("ctx"), j.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(
			j.Qual("context", "Background").Call(),
			j.Id("s").Dot("shutdownTimeout"),
		),
		j.Id("defer").Id("cancel").Call(),
		j.Line(),
		j.Return(j.Id("s").Dot("server").Dot("Shutdown").Call(j.Id("ctx"))),
	)

	return f
}
