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
	Path     string
	Module   string
	Adapters []string
}

func (s *Settings) IsAdapterChecked(adapterName string) bool {
	for _, adapter := range s.Adapters {
		if adapter == adapterName {
			return true
		}
	}

	return false
}
