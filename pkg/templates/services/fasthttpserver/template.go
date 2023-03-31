package fasthttpserver

import (
	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
)

type Service struct {
	name             string // name of the service
	*domain.Settings        // settings of the project
}

// GetName returns the name of the service
func (s *Service) GetName() string {
	return s.name
}

func NewService(name string, settings *domain.Settings) domain.ModuleI {
	return &Service{
		name:     name,
		Settings: settings,
	}
}

func (s *Service) ConfigGo() *j.Statement {
	return nil
}

func (s *Service) ConfigYAML() map[string]interface{} {
	return nil
}

func (s *Service) AppInit() []j.Code {
	return []j.Code{
		j.Line(),
		j.Id("httpServer").Op(":=").Qual(s.Module+"/pkg/fasthttpserver", "NewServer").Params(),
	}
}

func (s *Service) AppShutdown() []j.Code {
	return []j.Code{}
}

func (s *Service) Service() *j.File {
	f := j.NewFilePathName(s.Settings.Module+"/pkg/redis", "redis")

	// Service struct
	sStruct := j.Type().Id("Redis").Struct(
		j.Id("Client").Qual("github.com/go-redis/redis/v8", "Client"),
	)

	f.Add(sStruct)

	return f
}
