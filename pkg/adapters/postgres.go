package adapters

import (
	j "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type PostgresAdapter struct {
	name        string // name of the adapter
	displayName string // name of the adapter that will be displayed in the CLI
}

// GetName returns the name of the adapter
func (adp *PostgresAdapter) GetName() string {
	return adp.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (adp *PostgresAdapter) GetDisplayName() string {
	return adp.displayName
}

func NewPostgresAdapter() domain.ModuleI {
	return &PostgresAdapter{
		name:        "postgres",
		displayName: "Postgres",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (adp *PostgresAdapter) ConfigYAML() map[string]interface{} {
	return map[string]interface{}{
		"postgres": map[string]interface{}{
			"url":           "postgresql://user@localhost",
			"max_pool_size": 10,
		},
	}
}

// ConfigGo is the configuration of the adapter in Go format
func (adp *PostgresAdapter) ConfigGo() *j.Statement {
	return j.Id("Postgres").Struct(
		j.Id("URL").String().Tag(map[string]string{"mapstructure": "url", "json": "url"}),
		j.Id("MaxPoolSize").Int().Tag(map[string]string{"mapstructure": "max_pool_size", "json": "max_pool_size"}),
	).Tag(map[string]string{"mapstructure": "postgres", "json": "postgres"})
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (adp *PostgresAdapter) AppInit(module string) []j.Code {
	return []j.Code{
		j.List(j.Id("pg"), j.Err()).Op(":=").Qual(module+"/pkg/postgres", "New").Params(j.Id("gCtx"), j.Id("cfg.Postgres.URL")),
		j.Line(),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(j.Lit("app - Run - postgres.New"), j.Lit("error"), j.Id("err")),
		),
		j.Line(),
	}
}

// AppSelect - Each AppSelect branch is apart of a bigger switch statement that's in the internal/app/app.go Run() function
func (adp *PostgresAdapter) AppSelect(module string) j.Code {
	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (adp *PostgresAdapter) AppShutdown(module string) []j.Code {
	return []j.Code{
		j.Id("pg").Dot("Close").Call(),
	}
}

// Service is the code that will be added to its own `pkg` folder
func (adp *PostgresAdapter) Service(module string) *j.File {
	f := j.NewFilePathName(module+"/pkg/postgres", "postgres")

	// Service struct
	sStruct := j.Type().Id("Postgres").Struct(
		j.Id("maxPoolSize").Int(),
		j.Id("connAttempts").Int(),
		j.Id("connTimeout").Qual("time", "Duration"),
		j.Line(),
		j.Id("Pool").Add(utils.Jptr).Qual("github.com/jackc/pgx/v5/pgxpool", "Pool"),
	)

	f.Add(sStruct)

	f.Func().Id("New").Params(j.List(j.Id("ctx").Qual("context", "Context"), j.Id("url").String())).Params(j.Add(utils.Jptr).Id("Postgres"), j.Error()).Block(
		j.Id("pg").Op(":=").Add(utils.Rptr).Id("Postgres").Values(j.Dict{
			j.Id("maxPoolSize"):  j.Lit(10),
			j.Id("connAttempts"): j.Lit(10),
			j.Id("connTimeout"):  j.Qual("time", "Second").Op("*").Lit(5),
		}),
		j.Line(),
		j.List(j.Id("poolConfig"), j.Err()).Op(":=").Qual("github.com/jackc/pgx/v5/pgxpool", "ParseConfig").Params(j.Id("url")),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Err()),
		),
		j.Line(),
		j.Id("poolConfig").Dot("MaxConns").Op("=").Int32().Parens(j.Id("pg.maxPoolSize")),
		j.Line(),
		j.For(
			j.Id("pg").Dot("connAttempts").Op(">").Lit(0).Block(
				j.List(j.Id("pg.Pool"), j.Err()).Op("=").Qual("github.com/jackc/pgx/v5/pgxpool", "NewWithConfig").Call(j.Id("ctx"), j.Id("poolConfig")),
				j.If(j.Err().Op("==").Nil()).Block(
					j.Return(j.Id("pg"), j.Nil()),
				),
				j.Line(),
				j.Qual("time", "Sleep").Call(j.Id("pg").Dot("connTimeout")),
				j.Id("pg").Dot("connAttempts").Op("--"),
			),
		),
		j.Line(),
		j.If(j.Err().Op("!=").Nil()).Block(
			j.Return(j.Nil(), j.Err()),
		),
		j.Line(),
		j.Return(j.Id("pg"), j.Nil()),
	)

	f.Line()

	f.Func().Params(j.Id("pg").Op("*").Id("Postgres")).Id("Close").Call().Block(
		j.If(j.Id("pg.Pool").Op("!=").Nil()).Block(
			j.Id("pg.Pool").Dot("Close").Call(),
		),
	)

	return f
}
