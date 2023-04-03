package services

import (
	"github.com/mahcks/gowizard/pkg/domain"
	flavors "github.com/mahcks/gowizard/pkg/flavors/gql"
)

type GQLService struct {
	name        string // name of the service
	displayName string // name of the adapter that will be displayed in the CLI
	flavors     map[string]domain.FlavorI
}

func NewGQLService() domain.ServiceI {
	return &GQLService{
		name:        "gql",
		displayName: "GQL",
		flavors: map[string]domain.FlavorI{
			"gqlgen": flavors.NewGQLGenFlavor(),
		},
	}
}

// GetName returns the name of the service
func (svc *GQLService) GetName() string {
	return svc.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (svc *GQLService) GetDisplayName() string {
	return svc.displayName
}

// GetFlavors - returns the flavors that are available for this service
func (svc *GQLService) GetFlavors() map[string]domain.FlavorI {
	return svc.flavors
}

// GetFlavor - returns the flavor that is available for this service
func (svc *GQLService) GetFlavor(flavor string) domain.FlavorI {
	return svc.flavors[flavor]
}
