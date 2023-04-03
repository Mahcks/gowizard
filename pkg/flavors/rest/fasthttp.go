package services

import (
	"os"

	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type FastHTTPFlavor struct {
	name        string // name of the flavor
	displayName string // name of the adapter that will be displayed in the CLI
	description string // description of the flavor
}

// GetName returns the name of the flavor
func (flv *FastHTTPFlavor) GetName() string {
	return flv.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (flv *FastHTTPFlavor) GetDisplayName() string {
	return flv.displayName
}

// GetDescription - returns the description of the flavor
func (flv *FastHTTPFlavor) GetDescription() string {
	return flv.description
}

func NewFastHTTPFlavor() domain.FlavorI {
	return &FastHTTPFlavor{
		name:        "fasthttp",
		displayName: "valyala/fasthttp",
		description: "Fast HTTP package for Go. Tuned for high performance. Zero memory allocations in hot paths. Up to 10x faster than net/http",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (flv *FastHTTPFlavor) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the adapter in Go format
func (flv *FastHTTPFlavor) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (flv *FastHTTPFlavor) AppInit(module string) []j.Code {
	return []j.Code{
		j.Id("handler").Op(":=").Qual("github.com/fasthttp/router", "New").Call(),
		j.Line(),
		j.Id("httpServer").Op(":=").Qual(module+"/pkg/httpserver", "New").Call(j.Id("handler").Dot("Handler")),
	}
}

func (flv *FastHTTPFlavor) AppSelect(module string) j.Code {
	return j.Case(
		j.Id("err").Op("=").Op("<-").Id("httpServer").Dot("Notify").Call()).Block(
		j.Qual("fmt", "Println").Call(j.Lit("app.httpServer.Notifiy()"), j.Err()),
	)
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (flv *FastHTTPFlavor) AppShutdown(module string) []j.Code {
	return []j.Code{
		j.Id("err").Op("=").Id("httpServer").Dot("Shutdown").Call(),
		j.Line(),
		j.If(j.Id("err").Op("!=").Nil()).Block(
			j.Qual("fmt", "Println").Call(j.Lit("app.httpServer.Shutdown()"), j.Err()),
		),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (flv *FastHTTPFlavor) Service(module, path string) *j.File {
	f := j.NewFilePathName(module+"/pkg/httpserver", "httpserver")

	// Before saving the file, create the directories if they don't exist
	outputPath := path + "/pkg/httpserver"
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		utils.PrintError("error creating directories: %s", err)
		return nil
	}

	// Service struct
	sStruct := j.Type().Id("Service").Struct(
		j.Id("server").Add(utils.Jptr).Qual("github.com/valyala/fasthttp", "Server"),
		j.Id("notify").Chan().Error(),
		j.Id("shutdownTimeout").Qual("time", "Duration"),
	)

	f.Add(sStruct)

	f.Var().Id("defaultReadTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)
	f.Var().Id("defaultWriteTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)
	f.Var().Id("defaultAddr").Op("=").Lit("0.0.0.0:80")
	f.Var().Id("defaultShutdownTimeout").Op("=").Qual("time", "Second").Op("*").Lit(5)

	// New service
	f.Func().Id("New").Params(j.Id("handler").Qual("github.com/valyala/fasthttp", "RequestHandler")).Add(utils.Jptr).Id("Service").Block(
		j.Id("httpServer").Op(":=").Add(utils.Rptr).Qual("github.com/valyala/fasthttp", "Server").Values(j.Dict{
			j.Id("Handler"):      j.Id("handler"),
			j.Id("ReadTimeout"):  j.Id("defaultReadTimeout"),
			j.Id("WriteTimeout"): j.Id("defaultWriteTimeout"),
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
			j.Id("s").Dot("notify").Op("<-").Id("s").Dot("server").Dot("ListenAndServe").Call(j.Id("defaultAddr")),
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
		j.Return(j.Id("s").Dot("server").Dot("Shutdown").Call()),
	)

	err = f.Save(outputPath + "/server.go")
	if err != nil {
		utils.PrintError("error saving file: %s", err)
		return nil
	}

	return f
}
