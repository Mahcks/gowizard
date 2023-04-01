package domain

import (
	j "github.com/dave/jennifer/jen"
)

type ModuleI interface {
	// GetName returns the name of the module
	GetName() string

	// ConfigYAML is the configuration of the module in YAML format
	ConfigYAML() map[string]interface{}
	// ConfigGo is the configuration of the module in Go format
	ConfigGo() *j.Statement
	// AppInit is the code that will be added to the START internal/app/app.go Run() function
	AppInit() []j.Code
	// AppInit is the code that will be added to the END internal/app/app.go Run() function
	AppShutdown() []j.Code
	// Service is the code that will be added to its own `pkg` folder
	Service() *j.File
}

type Adapter struct {
	name string
	*Settings
}

type AdapterI interface {
	GetName() string
	GetSettings() *Settings
}

func (adp *Adapter) GetName() string {
	return adp.name
}

func (adp *Adapter) GetSettings() *Settings {
	return adp.Settings
}

type Settings struct {
	Path          string   // Path to the module
	Logger        string   // Logger name
	Module        string   // Module name
	ModuleVersion string   // Go module version
	Adapters      []string // Enabled adapters
	Services      []string // Enabled services
	Controllers   []string // Enabled controllers
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
	for _, service := range s.Services {
		if service == serviceName {
			return true
		}
	}

	return false
}
