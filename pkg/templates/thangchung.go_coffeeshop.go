package repos

import (
	"github.com/mahcks/gowizard/pkg/domain"
)

type GoCoffeshopRepo struct {
	name             string // name of the repo
	shortDescription string // short description of the repo (used in CLI selection)
}

func (r *GoCoffeshopRepo) GetName() string {
	return r.name
}

func (r *GoCoffeshopRepo) GetShortDescription() string {
	return r.shortDescription
}

func NewGoCoffeshopRepo() domain.TemplateI {
	return &GoCoffeshopRepo{
		name:             "github.com/thangchung/go-coffeeshop",
		shortDescription: "A practical event-driven microservices demo built with Golang. Nomad, Consul Connect, Vault, and Terraform for deployment",
	}
}

func (r *GoCoffeshopRepo) Setup(path string) error {
	return nil
}
