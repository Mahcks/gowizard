package cmd

import (
	"fmt"
	"os"

	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/mahcks/gowizard/pkg/ui"
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

		// Prompt for adapters
		adapters, err := ui.PromptForAdapters()
		if err != nil {
			return
		}

		// Propt for services
		services, err := ui.PromptForServices()
		if err != nil {
			return
		}

		// For each service selected, prompt for the service's adapters
		chosenFlavors := make(map[string]string, len(services)) // map[service]flavor
		for _, service := range services {
			flavor, err := ui.PromptForServiceFlavor(service)
			if err != nil {
				return
			}

			chosenFlavors[service] = flavor
		}

		if chosenFlavors == nil || adapters == nil {
			fmt.Println("No services or adapters selected.")
			return
		}

		gen.SetSettings(module, goVersion, path, adapters, chosenFlavors)

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
