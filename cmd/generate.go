package cmd

import (
	"fmt"
	"os"

	"github.com/mahcks/gowizard/pkg/generator"
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
			fmt.Println("Error: ", err)
			return
		}

		if moduleName == "" {
			fmt.Println("Error: module name is required.")
			return
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if path == "" {
			fmt.Println("Error: module path is required.")
			return
		}

		// Open the directory user has given
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// Get list of files/directories in the directory
		files, err := f.Readdir(-1)
		if err != nil {
			fmt.Println(err)
			return
		}

		// If directory is not empty, throw an error
		if len(files) != 0 {
			fmt.Println("Error: Directory is NOT empty")
		}

		// Fetch adapters from flags
		adapters, err := cmd.Flags().GetStringSlice("adapter")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		generator.NewGenerator(moduleName, path, adapters, []string{})
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("module", "m", "", "Name of the module")
	generateCmd.Flags().StringP("path", "p", "", "Path to the module")

	generateCmd.Flags().StringSliceP("adapter", "a", []string{}, "Add an adapter to the project, i.e. mariadb, redis")
}
