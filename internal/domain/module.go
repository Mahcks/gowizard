package domain

import (
	j "github.com/dave/jennifer/jen"
)

type ModuleI interface {
	ConfigGo() *j.Statement
	AppInit() []j.Code
	AppShutdown() *j.Statement
}

type Settings struct {
	Folder          string
	ProjectName     string
	EnabledAdapters []string // List of enabled adapters by user
}

func (s *Settings) IsAdapterChecked(adapterName string) bool {
	for _, adapter := range s.EnabledAdapters {
		if adapter == adapterName {
			return true
		}
	}

	return false
}
