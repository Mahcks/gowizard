package logger

import (
	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
)

type ZapLogger struct {
	name        string // name of the logger
	displayName string // name of the adapter that will be displayed in the CLI
}

// GetName returns the name of the logger
func (l *ZapLogger) GetName() string {
	return l.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (l *ZapLogger) GetDisplayName() string {
	return l.displayName
}

func NewZapLogger() domain.ModuleI {
	return &ZapLogger{
		name:        "zap",
		displayName: "Zap",
	}
}

// ConfigYAML is the configuration of the logger in YAML format
func (m *ZapLogger) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the logger in Go format
func (l *ZapLogger) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (m *ZapLogger) AppInit(module string) []j.Code {
	return nil
}

func (m *ZapLogger) AppSelect(module string) j.Code {

	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (m *ZapLogger) AppShutdown(module string) []j.Code {
	return nil
}

// Service is the code that will be added to its own `pkg` folder
func (a *ZapLogger) Service(module string) *j.File {
	f := j.NewFilePathName(module+"/pkg/logger", "zap")

	f.Func().Id("New").Params(j.Id("level").String()).Error().Block(
		j.Qual("log", "SetOutput").Call(j.Qual("io", "Discard")),
		j.Line(),
		j.Var().Id("lvl").Qual("go.uber.org/zap/zapcore", "Level"),
		j.Line(),
		j.Switch(j.Id("level")).Block(
			j.Case(j.Lit("debug")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "DebugLevel"),
			),
			j.Case(j.Lit("info")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "InfoLevel"),
			),
			j.Case(j.Lit("warn")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "WarnLevel"),
			),
			j.Case(j.Lit("error")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "ErrorLevel"),
			),
			j.Case(j.Lit("panic")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "PanicLevel"),
			),
			j.Case(j.Lit("fatal")).Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "FatalLevel"),
			),
			j.Default().Block(
				j.Id("lvl").Op("=").Qual("go.uber.org/zap", "InfoLevel"),
			),
		),
		j.Line(),
		j.Id("cfg").Op(":=").Qual("go.uber.org/zap", "NewProductionConfig").Call(),
		j.Id("cfg").Dot("Level").Op("=").Qual("go.uber.org/zap", "NewAtomicLevelAt").Call(j.Id("lvl")),
		j.Id("logger").Op(",").Id("err").Op(":=").Id("cfg").Dot("Build").Call(),
		j.If(j.Id("err").Op("!=").Nil()).Block(
			j.Return(j.Id("err")),
		),
		j.Line(),
		j.Qual("go.uber.org/zap", "ReplaceGlobals").Call(j.Id("logger")),
		j.Return(j.Nil()),
	)

	return f
}
