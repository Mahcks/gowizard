package services

import (
	"os"

	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type BeegoFlavor struct {
	name        string // name of the flavor
	displayName string // name of the adapter that will be displayed in the CLI
	description string // description of the flavor
}

// GetName returns the name of the flavor
func (flv *BeegoFlavor) GetName() string {
	return flv.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (flv *BeegoFlavor) GetDisplayName() string {
	return flv.displayName
}

// GetDescription - returns the description of the flavor
func (flv *BeegoFlavor) GetDescription() string {
	return flv.description
}

func NewBeegoFlavor() domain.FlavorI {
	return &BeegoFlavor{
		name:        "beego",
		displayName: "beego/beego",
		description: "beego is an open-source, high-performance web framework for the Go programming language.",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (flv *BeegoFlavor) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the adapter in Go format
func (flv *BeegoFlavor) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (flv *BeegoFlavor) AppInit(module string) []j.Code {
	return nil
}

func (flv *BeegoFlavor) AppSelect(module string) j.Code {
	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (flv *BeegoFlavor) AppShutdown(module string) []j.Code {
	return nil
}

// Service is the code that will be added to its own `pkg` folder
func (flv *BeegoFlavor) Service(module, path string) *j.File {
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
