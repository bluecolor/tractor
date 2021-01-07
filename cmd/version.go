package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "undefined"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Tractor",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tractor Data Ingestion Utility " + version)
	},
}
