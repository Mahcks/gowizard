package domain

import (
	j "github.com/dave/jennifer/jen"
)

type ModuleI interface {
	// GetName returns the name of the module
	GetName() string
	// GetDisplayName - what will be displayed in the CLI when prompted
	GetDisplayName() string

	// ConfigYAML is the configuration of the module in YAML format
	ConfigYAML() map[string]interface{}
	// ConfigGo is the configuration of the module in Go format
	ConfigGo() *j.Statement
	// AppInit is the code that will be added to the START internal/app/app.go Run() function
	AppInit(module string) []j.Code
	// AppSelect - Each AppSelect branch is apart of a bigger switch statement that's in the internal/app/app.go Run() function
	AppSelect(module string) j.Code
	// AppInit is the code that will be added to the END internal/app/app.go Run() function
	AppShutdown(module string) []j.Code
	// Service is the code that will be added to its own `pkg` folder
	Service(module, path string) *j.File
}

type ServiceI interface {
	// GetName returns the name of the module
	GetName() string
	// GetDisplayName - what will be displayed in the CLI when prompted
	GetDisplayName() string
	// GetFlavors - returns the flavors that are available for this service
	GetFlavors() map[string]FlavorI
	// GetFlavor - returns flavor by name
	GetFlavor(flavor string) FlavorI
}

type FlavorI interface {
	// GetName returns the name of the module
	GetName() string
	// GetDisplayName - what will be displayed in the CLI when prompted
	GetDisplayName() string
	// GetDescription - returns the description of the flavor
	GetDescription() string

	// ConfigYAML is the configuration of the module in YAML format
	ConfigYAML() map[string]interface{}
	// ConfigGo is the configuration of the module in Go format
	ConfigGo() *j.Statement
	// AppInit is the code that will be added to the START internal/app/app.go Run() function
	AppInit(module string) []j.Code
	// AppSelect - Each AppSelect branch is apart of a bigger switch statement that's in the internal/app/app.go Run() function
	AppSelect(module string) j.Code
	// AppInit is the code that will be added to the END internal/app/app.go Run() function
	AppShutdown(module string) []j.Code
	// Service is the code that will be added to its own `pkg` folder
	Service(module, path string) *j.File
}

type Settings struct {
	Path          string            // Path to the module
	Module        string            // Module name
	ModuleVersion string            // Go module version
	Adapters      []string          // Enabled adapters
	Services      map[string]string // Enabled services, key is the service name, value is the flavor name
	Controllers   []string          // Enabled controllers
}

// IsAdapterChecked checks if the adapter is enabled
func (s *Settings) IsAdapterChecked(adapterName string) bool {
	for _, adapter := range s.Adapters {
		if adapter == adapterName {
			return true
		}
	}

	return false
}

// IsServiceChecked checks if the service is enabled
func (s *Settings) IsServiceChecked(serviceName string) bool {
	for service := range s.Services {
		if service == serviceName {
			return true
		}
	}

	return false
}
