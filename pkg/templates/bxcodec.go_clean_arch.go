package repos

import (
	"github.com/mahcks/gowizard/pkg/domain"
)

type GoCleanArchTemplateRepo struct {
	name             string // name of the repo
	shortDescription string // short description of the repo (used in CLI selection)
}

func (r *GoCleanArchTemplateRepo) GetName() string {
	return r.name
}

func (r *GoCleanArchTemplateRepo) GetShortDescription() string {
	return r.shortDescription
}

func NewGoCleanArchTemplateRepo() domain.TemplateI {
	return &GoCleanArchTemplateRepo{
		name:             "github.com/bxcodec/go-clean-arch",
		shortDescription: "Go (Golang) Clean Architecture based on Reading Uncle Bob's Clean Architecture",
	}
}

func (r *GoCleanArchTemplateRepo) Setup(path string) error {
	return nil
}
