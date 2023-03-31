package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahcks/gowizard/pkg/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gowizard",
	Short: "A brief description of your application",
	Long: `gowizard is a CLI tool to generate Go modules with a setup wizard.

You can also just skip the wizard... 
gowizard generate --module github.com/user/repo --path /some/path --adapter mariadb,redis --service rest-fasthttp`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var enabledAdapters []string
		var enabledServices []string
		enabledAdapters = append(enabledAdapters, "mariadb", "redis")
		// enabledServices = append(enabledServices, "rest-fasthttp")

		// init styles; optional, just showing as a way to organize styles
		// start bubble tea and init first model
		questions := []Question{newShortQuestion("What is the name of your module?"), newShortQuestion("Enter your desired path for the module")}
		main := New(questions)

		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
			return
		}
		defer f.Close()

		p := tea.NewProgram(*main, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
			return
		}

		gen := generator.NewGenerator(questions[0].answer, questions[1].answer, enabledAdapters, enabledServices)
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
