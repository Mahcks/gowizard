package cmd

import (
	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/ui"
	"github.com/mahcks/gowizard/pkg/utils"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Use a predefined template to generate a project.",
	Run: func(cmd *cobra.Command, args []string) {
		gen := generator.NewGenerator()
		ui := ui.NewUI(gen)

		// Ask for module name
		module, err := ui.PromptForModuleName()
		if err != nil {
			return
		}

		// Ask for what version of Go to use
		goVersion, err := ui.PromptForGoVersion()
		if err != nil {
			return
		}

		// Ask for module path
		path, err := ui.PromptForModulePath()
		if err != nil {
			return
		}

		// Check if user wants to use a custom template
		isCustom, err := cmd.Flags().GetBool("custom")
		if err != nil {
			utils.PrintError("error getting custom flag: %s", err)
			return
		}

		var template string
		if isCustom {
			template, err = ui.PromptForCustomTemplate()
			if err != nil {
				return
			}
		} else {
			// Prompt for hard-coded template
			template, err = ui.PromptForTemplate()
			if err != nil {
				return
			}
		}

		gen.SetSettings(module, goVersion, path, []string{}, nil)
		err = gen.UseTemplate(template, isCustom)
		if err != nil {
			utils.PrintError("error generating template: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().BoolP("custom", "c", false, "Use a custom template, this will let to specify a Go version and module name but won't have any extra setup.")
}
