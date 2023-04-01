package cmd

import (
	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/utils"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new project.",
	Long:  `Generate a new Go module with a given name and path. You can also specifiy services and adapters to be included in the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, err := cmd.Flags().GetString("module")
		if err != nil {
			utils.PrintError("error getting module flag: %s", err)
			return
		}

		if moduleName == "" {
			utils.PrintError("module name is required")
			return
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			utils.PrintError("error getting path flag: %s", err)
			return
		}

		if path == "" {
			utils.PrintError("module path is required")
			return
		}

		// Open the directory user has given
		isEmpty, err := utils.IsDirEmpty(path)
		if err != nil {
			utils.PrintError("unable to open file: %s" + err.Error())
			return
		}

		if !isEmpty {
			utils.PrintError("The directory you have specified is not empty.")
			return
		}

		// Fetch adapters from flags
		adapters, err := cmd.Flags().GetStringSlice("adapter")
		if err != nil {
			utils.PrintError("error getting adapter flags: %s", err)
			return
		}

		// Get the version of Go to use, defaults to the users latest installed version
		goVersion, err := cmd.Flags().GetString("go-version")
		if err != nil {
			utils.PrintError("error getting go-version flag: %s", err)
			return
		}

		// Get the template to use
		template, err := cmd.Flags().GetString("template")
		if err != nil {
			utils.PrintError("error getting template flag: %s", err)
			return
		}

		gen := generator.NewGenerator(moduleName, goVersion, path, adapters, []string{})

		// If a template is specified, use it
		if template != "" {
			err = gen.UseTemplate(template)
			if err != nil {
				utils.PrintError("error setting template: %s", err)
				return
			}

			return
		}

		err = gen.Generate()
		if err != nil {
			utils.PrintError("%s", err)

			errRollback := gen.Rollback()
			if errRollback != nil {
				utils.PrintError("error rolling back: %s", errRollback)
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Get the current version of Go on the users system
	cmdVersion, err := utils.GetGoVersion()
	if err != nil {
		utils.PrintError("error getting Go version: %s", err)
		return
	}

	generateCmd.Flags().StringP("module", "m", "", "Name of the module")
	generateCmd.Flags().StringP("path", "p", "./", "Path to the module")
	generateCmd.Flags().StringP("go-version", "v", cmdVersion, "Go version to use - defaults to your latest installed version")
	generateCmd.Flags().StringP("template", "t", "", "Template to use for the project")

	generateCmd.Flags().StringSliceP("adapter", "a", []string{}, "Add an adapter to the project, i.e. mariadb, redis")
}
