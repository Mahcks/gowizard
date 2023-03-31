package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/utils"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gowizard",
	Short: "A brief description of your application",
	Long: `gowizard is a CLI tool to generate Go modules with a setup wizard.

You can also just skip the wizard... 
gowizard generate --module github.com/user/repo --path /some/path --adapter mariadb,redis --service rest-fasthttp`,
	Run: func(cmd *cobra.Command, args []string) {
		iconStyles := survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "[?]"
			icons.Question.Format = "magenta+b"

			icons.MarkedOption.Format = "cyan+b"
		})

		// Ask for module name
		module := ""
		prompt := &survey.Input{
			Message: "What is your desired module name?",
			Help:    "This is the name of the module that will be generated. It should be in the format of github.com/userororg/repo",
		}
		err := survey.AskOne(prompt, &module, iconStyles, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Ask for module path
		path := ""
		prompt = &survey.Input{
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
		err = survey.AskOne(prompt, &path, iconStyles, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Ensure that the directory is empty before generating
		isDirEmpty, err := utils.IsDirEmpty(path)

		if !isDirEmpty {
			fmt.Println(ansi.Color(fmt.Sprintf(`[✗] Error: directory "%s" is not empty`, path), "red"), ansi.ColorCode("reset"))
			return
		}

		// Prompt for adapters
		adapters := []string{}
		adapterPrompt := &survey.MultiSelect{
			Message: "Choose adapters:",
			Options: []string{"MariaDB", "MongoDB", "Redis"},
		}
		err = survey.AskOne(adapterPrompt, &adapters, iconStyles)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Lowercase all adapter names before passing them to the generator
		for i, adapter := range adapters {
			adapters[i] = strings.ToLower(adapter)
		}

		gen := generator.NewGenerator(module, path, adapters, []string{})
		err = gen.Generate()
		if err != nil {
			fmt.Println(err.Error())
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gowizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
