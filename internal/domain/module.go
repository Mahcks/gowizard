package domain

import (
	j "github.com/dave/jennifer/jen"
)

type ModuleI interface {
	ConfigGo() *j.Statement
	AppInit() []j.Code
	AppShutdown() *j.Statement
}
