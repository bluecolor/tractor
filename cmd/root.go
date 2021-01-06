package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tractor",
	Short: "Tractor is an easy to use data ingestion tool",
	Long: `A Fast and Flexible Data Ingestion Utility built with
love by blue and friends in Go.
Complete documentation is available at http://github.io/bluecolor/tractor`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
