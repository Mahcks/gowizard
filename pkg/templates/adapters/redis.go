package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type RedisAdapter struct {
	name string // name of the adapter
}

// GetName returns the name of the adapter
func (a *RedisAdapter) GetName() string {
	return a.name
}

func NewRedisAdapter() domain.ModuleI {
	return &RedisAdapter{
		name: "redis",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (a *RedisAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"redis": map[string]interface{}{
			"host":     "localhost",
			"port":     "6379",
			"password": "password123",
		},
	}
}

// ConfigGo is the configuration of the adapter in Go format
func (a *RedisAdapter) ConfigGo() *j.Statement {
	return j.Id("Redis").Struct(
		j.Id("Host").String().Tag(map[string]string{"mapstructure": "host", "json": "host"}),
		j.Id("Port").String().Tag(map[string]string{"mapstructure": "port", "json": "port"}),
		j.Id("Password").String().Tag(map[string]string{"mapstructure": "password", "json": "password"}),
	).Tag(map[string]string{"mapstructure": "redis", "json": "redis"})
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (a *RedisAdapter) AppInit(module string) []j.Code {
	return []j.Code{
		j.Line(),
		j.List(j.Id("redisClient"), j.Err()).Op(":=").Qual(module+"/pkg/redis", "New").Params(j.Id("gCtx"), j.Id("cfg.Redis.Host"), j.Id("cfg.Redis.Port"), j.Id("cfg.Redis.Password")).Op(";"),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - redis.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
		j.Line(),
		j.Qual("go.uber.org/zap", "S").Call().Dot("Infow").Call(j.Lit("main - app - Run"), j.Lit("message"), j.Lit("Redis initialized")),
		j.Line(),
	}
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (a *RedisAdapter) AppShutdown() []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("redisClient").Dot("Close").Call(),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (a *RedisAdapter) Service(module string) *j.File {
	f := j.NewFilePathName(module+"/pkg/redis", "redis")

	// Service struct
	sStruct := j.Type().Id("Redis").Struct(
		j.Id("Client").Add(utils.Jptr).Qual("github.com/go-redis/redis/v8", "Client"),
	)

	f.Add(sStruct)

	// New function
	f.Func().Id("New").Params(
		j.Id("ctx").Qual("context", "Context"),
		j.Id("host"),
		j.Id("port"),
		j.Id("password").String(),
	).Op("(").List(j.Op("*").Add(j.Id("Redis"), j.Op(","), j.Error()).Op(")")).Block(
		j.Id("client").Op(":=").Qual("github.com/go-redis/redis/v8", "NewClient").Params(
			j.Op("&").Qual("github.com/go-redis/redis/v8", "Options").Values(
				j.Dict{
					j.Id("Addr"):     j.Id("host").Op("+").Lit(":").Op("+").Id("port"),
					j.Id("Password"): j.Id("password"),
					j.Id("DB"):       j.Lit(0),
				},
			),
		),
		j.Line(),
		j.List(j.Id("_"), j.Id("err")).Op(":=").Id("client").Dot("Ping").Call(j.Id("ctx")).Dot("Result").Call(),
		j.If(j.Id("err").Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Id("err")),
		),
		j.Line(),
		j.Return(
			j.Op("&").Id("Redis").Values(
				j.Dict{
					j.Id("Client"): j.Id("client"),
				},
			),
			j.Nil(),
		),
	)

	f.Add(j.Line())

	// Close function
	f.Func().Params(j.Id("r").Op("*").Id("Redis")).Id("Close").Params().Error().Block(
		j.If(
			j.Id("r").Dot("Client").Op("!=").Nil().Block(
				j.Id("err").Op(":=").Id("r").Dot("Client").Dot("Close").Call(),
				j.If(j.Id("err").Op("!=").Nil()).Block(
					j.Return(j.Id("err")),
				),
			),
		),
		j.Return(j.Nil()),
	)

	return f
}
