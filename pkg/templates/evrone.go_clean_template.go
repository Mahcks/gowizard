package repos

import (
	"github.com/mahcks/gowizard/pkg/domain"
)

type GoCleanTemplateRepo struct {
	name             string // name of the repo
	shortDescription string // short description of the repo (used in CLI selection)
}

func (r *GoCleanTemplateRepo) GetName() string {
	return r.name
}

func (r *GoCleanTemplateRepo) GetShortDescription() string {
	return r.shortDescription
}

func NewGoCleanTemplateRepo() domain.TemplateI {
	return &GoCleanTemplateRepo{
		name:             "github.com/evrone/go-clean-template",
		shortDescription: "Clean Architecture template for Golang services",
	}
}

func (r *GoCleanTemplateRepo) Setup(path string) error {
	return nil
}
