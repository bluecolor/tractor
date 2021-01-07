package cmd

import (
	"github.com/bluecolor/tractor/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.AddCommand(mappingsCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints the info requested with argumet",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var mappingsCmd = &cobra.Command{
	Use:   "mappings",
	Short: "print mapping names to terminal",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := "mappings.yml"
		util.PrintMappingList(configFile)
	},
}
