package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gowizard",
	Long:  `All software has versions. This is gowizards's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gowizard: v" + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
