package redis

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/internal/domain"
)

type Adapter struct {
	domain.Settings
}

func NewAdapter(settings *domain.Settings) domain.ModuleI {
	return &Adapter{
		Settings: *settings,
	}
}

func (a *Adapter) ConfigGo() *j.Statement {
	return j.Id("Redis").Struct(
		j.Id("Host").String().Tag(map[string]string{"mapstructure": "host", "json": "host"}),
		j.Id("Port").String().Tag(map[string]string{"mapstructure": "port", "json": "port"}),
		j.Id("Password").String().Tag(map[string]string{"mapstructure": "password", "json": "password"}),
	).Tag(map[string]string{"mapstructure": "redis", "json": "redis"})
}

func (a *Adapter) AppInit() []j.Code {
	return []j.Code{
		j.List(j.Id("redisClient"), j.Err()).Op(":=").Qual(a.ProjectName+"/pkg/redis", "New").Params(j.Id("gCtx"), j.Id("cfg.Redis.Host"), j.Id("cfg.Redis.Port"), j.Id("cfg.Redis.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - redis.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("Redis initialized")),
	}
}

func (a *Adapter) AppShutdown() *j.Statement {
	return j.Id("redisClient").Dot("Close").Call()
}
