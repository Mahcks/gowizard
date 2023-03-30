/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mahcks/gowizard/internal/builder"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new project.",
	Long:  `Generate a new Go module with a given name and path. You can also specifiy services and adapters to be included in the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if projectName == "" {
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

		// Store enabled adapters
		var enabledAdapters []string

		mariaDBEnabled, err := cmd.Flags().GetBool("mariadb")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		redisAdapterEnabled, err := cmd.Flags().GetBool("redis")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if mariaDBEnabled {
			enabledAdapters = append(enabledAdapters, "mariadb")
		}

		if redisAdapterEnabled {
			enabledAdapters = append(enabledAdapters, "redis")
		}

		enabledAdapters = append(enabledAdapters, "logger")
		builder.NewBuilder(path, projectName, enabledAdapters)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().String("name", "", "Name of the module")
	generateCmd.Flags().String("path", "", "Path to the module")

	generateCmd.Flags().BoolP("mariadb", "", false, "Include MariaDB adapter")
	generateCmd.Flags().BoolP("redis", "", false, "Include Redis adapter")
}
