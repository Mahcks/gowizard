package repos

import (
	"github.com/mahcks/gowizard/pkg/domain"
)

type GoBackendCleanArchitectureTemplateRepo struct {
	name             string // name of the repo
	shortDescription string // short description of the repo (used in CLI selection)
}

func (r *GoBackendCleanArchitectureTemplateRepo) GetName() string {
	return r.name
}

func (r *GoBackendCleanArchitectureTemplateRepo) GetShortDescription() string {
	return r.shortDescription
}

func NewGoBackendCleanArchitectureTemplateRepo() domain.TemplateI {
	return &GoBackendCleanArchitectureTemplateRepo{
		name:             "github.com/amitshekhariitbhu/go-backend-clean-architecture",
		shortDescription: "A Go (Golang) Backend Clean Architecture project with Gin, MongoDB, JWT Authentication Middleware, Test, and Docker.",
	}
}

func (r *GoBackendCleanArchitectureTemplateRepo) Setup(path string) error {
	return nil
}
