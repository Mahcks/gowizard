package logger

import (
	"fmt"

	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
)

type ZapLogger struct {
	name             string // name of the logger
	*domain.Settings        // settings of the project
}

func (l *ZapLogger) GetName() string {
	return l.name
}

func NewZapLogger(settings *domain.Settings) domain.ModuleI {
	return &ZapLogger{
		name:     "zap",
		Settings: settings,
	}
}

func (l *ZapLogger) ConfigGo() *j.Statement {
	return nil
}

func (m *ZapLogger) ConfigYAML() map[string]interface{} {
	return nil
}

func (m *ZapLogger) AppInit() []j.Code {
	return nil
}

func (m *ZapLogger) AppShutdown() []j.Code {
	return nil
}

func (a *ZapLogger) Service() *j.File {
	f := j.NewFilePathName(a.Settings.Module+"/pkg/logger", "zap")

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

	err := f.Save(a.Settings.Path + "/pkg/logger/zap.go")
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
