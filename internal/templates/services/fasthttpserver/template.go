package fasthttpserver

import (
	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/internal/domain"
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
