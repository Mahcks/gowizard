package services

import (
	"github.com/mahcks/gowizard/pkg/domain"
	flavors "github.com/mahcks/gowizard/pkg/flavors/rest"
)

type RESTService struct {
	name        string // name of the service
	displayName string // name of the adapter that will be displayed in the CLI
	flavors     map[string]domain.FlavorI
}

func NewRESTService() domain.ServiceI {
	return &RESTService{
		name:        "rest",
		displayName: "REST",
		flavors: map[string]domain.FlavorI{
			"beego":    flavors.NewBeegoFlavor(),
			"fasthttp": flavors.NewFastHTTPFlavor(),
			"fiber":    flavors.NewFiberFlavor(),
			"gin":      flavors.NewGinFlavor(),
		},
	}
}

// GetName returns the name of the service
func (svc *RESTService) GetName() string {
	return svc.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (svc *RESTService) GetDisplayName() string {
	return svc.displayName
}

// GetFlavors - returns the flavors that are available for this service
func (svc *RESTService) GetFlavors() map[string]domain.FlavorI {
	return svc.flavors
}

// GetFlavor - returns the flavor that is available for this service
func (svc *RESTService) GetFlavor(flavor string) domain.FlavorI {
	return svc.flavors[flavor]
}
