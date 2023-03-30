package domain

import (
	j "github.com/dave/jennifer/jen"
)

type ModuleI interface {
	GetName() string
	ConfigGo() *j.Statement
	AppInit() []j.Code
	AppShutdown() []j.Code
}

type Settings struct {
	Path        string   // Path to the module
	Module      string   // Module name
	Adapters    []string // Enabled adapters
	Services    []string // Enabled services
	Controllers []string // Enabled controllers
}

func (s *Settings) IsAdapterChecked(adapterName string) bool {
	for _, adapter := range s.Adapters {
		if adapter == adapterName {
			return true
		}
	}

	return false
}

func (s *Settings) IsServiceChecked(serviceName string) bool {
	for _, service := range s.Services {
		if service == serviceName {
			return true
		}
	}

	return false
}
