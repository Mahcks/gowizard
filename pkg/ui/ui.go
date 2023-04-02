package ui

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/utils"
)

type UI struct {
	iconStyles survey.AskOpt
	gen        *generator.Generator
}

func NewUI(gen *generator.Generator) *UI {
	return &UI{
		iconStyles: survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "[?]"
			icons.Question.Format = "magenta+b"

			icons.MarkedOption.Format = "cyan+b"
		}),
		gen: gen,
	}
}

// PromptForModuleName prompts the user for the module name
func (ui *UI) PromptForModuleName() (string, error) {
	module := ""
	promptModule := &survey.Input{
		Message: "What is your desired module name?",
		Help:    "This is the name of the module that will be generated. It should be in the format of github.com/userororg/repo",
		Default: "github.com/mahcks/somemodule",
	}
	err := survey.AskOne(promptModule, &module, ui.iconStyles, survey.WithValidator(survey.Required))
	if err != nil {
		utils.PrintError("error prompting for module name: %s", err)

		return "", fmt.Errorf("error prompting for module name: %s", err)
	}

	return module, nil
}

// PromptForGoVersion prompts the user for the Go version to use
func (ui *UI) PromptForGoVersion() (string, error) {
	// Get the current version of Go on the users system
	cmdVersion, err := utils.GetGoVersion()
	if err != nil {
		utils.PrintError("error getting Go version: %s", err)

		return "", fmt.Errorf("error getting Go version: %s", err)
	}

	goVersion := ""
	promptGoVersion := &survey.Input{
		Message: "What version of Go would you like to use?",
		Default: cmdVersion,
	}
	err = survey.AskOne(promptGoVersion, &goVersion, ui.iconStyles, survey.WithValidator(survey.Required))
	if err != nil {
		utils.PrintError("error prompting for Go version: %s", err)

		return "", fmt.Errorf("error prompting for Go version: %s", err)
	}

	return goVersion, nil
}

// PromptForModulePath prompts the user for the module path
func (ui *UI) PromptForModulePath() (string, error) {
	path := ""
	promptPath := &survey.Input{
		Message: "Where would you like to place the module?",
		Default: "./",
		Suggest: func(toComplete string) []string {
			// Suggest directories in the current working directory
			files, err := os.ReadDir(".")
			if err != nil {
				return nil
			}

			var suggestions []string
			for _, file := range files {
				if file.IsDir() {
					suggestions = append(suggestions, file.Name())
				}
			}

			return suggestions
		}}
	err := survey.AskOne(promptPath, &path, ui.iconStyles, survey.WithValidator(survey.Required))
	if err != nil {
		utils.PrintError("error prompting for module path: %s", err)
		return "", fmt.Errorf("error prompting for module path: %s", err)
	}

	// Ensure that the directory is empty before generating
	isDirEmpty, err := utils.IsDirEmpty(path)

	if !isDirEmpty {
		return "", fmt.Errorf("directory %s is not empty", path)
	}

	return path, nil
}

func (ui *UI) PromptForAdapters() ([]string, error) {
	adapters := []string{}

	var options []string
	descriptions := make(map[string]string) // Use a map to store the descriptions
	for key, value := range ui.gen.GetAdapters() {
		options = append(options, key)
		descriptions[key] = value.GetName() // Add the description for each template
	}

	// Sort the options slice in alphabetical order
	sort.Strings(options)

	adapterPrompt := &survey.MultiSelect{
		Message: "Choose adapters:",
		Options: options,
	}
	err := survey.AskOne(adapterPrompt, &adapters, ui.iconStyles)
	if err != nil {
		utils.PrintError("error prompting for adapters: %s", err)
		return nil, fmt.Errorf("error prompting for adapters: %s", err)
	}

	// Lowercase all adapter names before passing them to the generator
	for i, adapter := range adapters {
		adapters[i] = strings.ToLower(adapter)
	}

	return adapters, nil
}

// PromptForTemplate prompts the user for the template to use that's in the repos template directory
func (ui *UI) PromptForTemplate() (string, error) {
	var options []string
	descriptions := make(map[string]string) // Use a map to store the descriptions
	for key, value := range ui.gen.GetTemplates() {
		options = append(options, key)
		descriptions[key] = value.GetShortDescription() // Add the description for each template
	}

	// Sort the options slice in alphabetical order
	sort.Strings(options)

	template := ""
	prompt := &survey.Select{
		Message: "Select a template:",
		Options: options,
		Description: func(value string, index int) string {
			return descriptions[value]
		},
	}
	survey.AskOne(prompt, &template, ui.iconStyles)

	return template, nil
}

// PromptForCustomTemplate prompts the user for the repostiory to setup
func (ui *UI) PromptForCustomTemplate() (string, error) {
	customTemplate := ""
	promptModule := &survey.Input{
		Message: "Enter the path to your custom template:",
		Help:    "Path to the Github repo that contains your custom template. i.e. https://github.com/user/template",
	}
	err := survey.AskOne(promptModule, &customTemplate, ui.iconStyles, survey.WithValidator(survey.Required))
	if err != nil {
		utils.PrintError("error prompting for custom template: %s", err)

		return "", fmt.Errorf("error prompting for custom template repository: %s", err)
	}

	// Parse the repository URL and remove https:// and .git if present since it's handled later
	parsedURL, err := url.Parse(customTemplate)
	if err != nil {
		return "", err
	}

	repoPath := strings.TrimSuffix(parsedURL.Host+parsedURL.Path, ".git")

	return repoPath, nil
}
