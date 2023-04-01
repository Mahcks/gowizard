package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	Version         = "0.1.0"
	versionTemplate = `gowizard: v{{.Version}}
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gowizard",
	Short: "A brief description of your application",
	Long: `gowizard is a CLI tool to generate Go modules with a setup wizard.

You can also just skip the wizard... 
gowizard generate --module github.com/username/module --path /path/to/module --adapter mariadb,redis,mongodb`,
	Run: func(cmd *cobra.Command, args []string) {

		// Ask for module name
		module := ""
		promptModule := &survey.Input{
			Message: "What is your desired module name?",
			Help:    "This is the name of the module that will be generated. It should be in the format of github.com/userororg/repo",
		}
		err := survey.AskOne(promptModule, &module, utils.IconStyles, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Get the current version of Go on the users system
		cmdVersion, err := utils.GetGoVersion()
		if err != nil {
			utils.PrintError("error getting Go version: %s", err)
			return
		}

		// Ask for what version of Go to use
		goVersion := ""
		promptGoVersion := &survey.Input{
			Message: "What version of Go would you like to use?",
			Default: cmdVersion,
		}
		err = survey.AskOne(promptGoVersion, &goVersion, utils.IconStyles, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Ask for module path
		path := ""
		promptPath := &survey.Input{
			Message: "Where would you like to place the module?",
			Default: "./",
			Suggest: func(toComplete string) []string {
				// Suggest directories in the current working directory
				files, err := os.ReadDir(".")
				if err != nil {
					fmt.Println(err)
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
		err = survey.AskOne(promptPath, &path, utils.IconStyles, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Ensure that the directory is empty before generating
		isDirEmpty, err := utils.IsDirEmpty(path)

		if !isDirEmpty {
			utils.PrintError(`directory "%s" is not empty`, path)
			return
		}

		// Ask if the user wants to use a template
		useTemplate := false
		prompt := &survey.Confirm{
			Message: "Use a pre-built template?",
		}
		survey.AskOne(prompt, &useTemplate, utils.IconStyles)

		// If the user wants to use a template, prompt for which one
		if useTemplate {
			gen := generator.NewGenerator(module, goVersion, path, []string{}, []string{})
			templates := gen.GetTemplates()

			var options []string
			descriptions := make(map[string]string) // Use a map to store the descriptions
			for key, value := range templates {
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
			survey.AskOne(prompt, &template, utils.IconStyles)

			err = gen.UseTemplate(template)
			if err != nil {
				utils.PrintError("error using template: %s", err)
				return
			}

			return
		}

		// Prompt for adapters
		adapters := []string{}
		adapterPrompt := &survey.MultiSelect{
			Message: "Choose adapters:",
			Options: []string{"MariaDB", "MongoDB", "Redis"},
		}
		err = survey.AskOne(adapterPrompt, &adapters, utils.IconStyles)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Lowercase all adapter names before passing them to the generator
		for i, adapter := range adapters {
			adapters[i] = strings.ToLower(adapter)
		}

		gen := generator.NewGenerator(module, goVersion, path, adapters, []string{})
		err = gen.Generate()
		if err != nil {
			fmt.Println(err.Error())

			errRollback := gen.Rollback()
			if errRollback != nil {
				utils.PrintError("error rolling back: %s", errRollback)
			}
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(versionTemplate)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gowizard.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gowizard" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gowizard")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
