package services

import (
	"os"

	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type FiberFlavor struct {
	name        string // name of the flavor
	displayName string // name of the adapter that will be displayed in the CLI
	description string // description of the flavor
}

// GetName returns the name of the flavor
func (flv *FiberFlavor) GetName() string {
	return flv.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (flv *FiberFlavor) GetDisplayName() string {
	return flv.displayName
}

// GetDescription - returns the description of the flavor
func (flv *FiberFlavor) GetDescription() string {
	return flv.description
}

func NewFiberFlavor() domain.FlavorI {
	return &FiberFlavor{
		name:        "fiber",
		displayName: "gofiber/fiber",
		description: "Express inspired web framework written in Go",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (flv *FiberFlavor) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the adapter in Go format
func (flv *FiberFlavor) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (flv *FiberFlavor) AppInit(module string) []j.Code {
	return nil
}

func (flv *FiberFlavor) AppSelect(module string) j.Code {
	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (flv *FiberFlavor) AppShutdown(module string) []j.Code {
	return nil
}

// Service is the code that will be added to its own `pkg` folder
func (flv *FiberFlavor) Service(module, path string) *j.File {
	f := j.NewFilePathName(module+"/pkg/httpserver", "httpserver")

	// Before saving the file, create the directories if they don't exist
	outputPath := path + "/pkg/httpserver"
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		utils.PrintError("error creating directories: %s", err)
		return nil
	}

	err = f.Save(outputPath + "/server.go")
	if err != nil {
		utils.PrintError("error saving file: %s", err)
		return nil
	}

	return f
}
