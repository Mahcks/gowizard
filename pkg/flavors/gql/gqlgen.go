package services

import (
	"os"

	j "github.com/dave/jennifer/jen"

	"github.com/mahcks/gowizard/pkg/domain"
	"github.com/mahcks/gowizard/pkg/utils"
)

type Gin struct {
	name        string // name of the flavor
	displayName string // name of the adapter that will be displayed in the CLI
	description string // description of the flavor
}

// GetName returns the name of the flavor
func (flv *Gin) GetName() string {
	return flv.name
}

// GetDisplayName - what will be displayed in the CLI when prompted
func (flv *Gin) GetDisplayName() string {
	return flv.displayName
}

// GetDescription - returns the description of the flavor
func (flv *Gin) GetDescription() string {
	return flv.description
}

func NewGQLGenFlavor() domain.FlavorI {
	return &Gin{
		name:        "gqlgen",
		displayName: "gqlgen",
	}
}

// ConfigYAML is the configuration of the adapter in YAML format
func (flv *Gin) ConfigYAML() map[string]interface{} {
	return nil
}

// ConfigGo is the configuration of the adapter in Go format
func (flv *Gin) ConfigGo() *j.Statement {
	return nil
}

// AppInit is the code that will be added to the START internal/app/app.go Run() function
func (flv *Gin) AppInit(module string) []j.Code {
	return nil
}

func (flv *Gin) AppSelect(module string) j.Code {
	return nil
}

// AppShutdown is the code that will be added to the END internal/app/app.go Run() function
func (flv *Gin) AppShutdown(module string) []j.Code {
	return nil
}

// Service is the code that will be added to its own `pkg` folder
func (flv *Gin) Service(module, path string) *j.File {
	f := j.NewFilePathName(module+"/pkg/gql", "gql")

	// Before saving the file, create the directories if they don't exist
	outputPath := path + "/pkg/gqlserver"
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
