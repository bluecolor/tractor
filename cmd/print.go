package cmd

import (
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.AddCommand(printMappingsCmd)
	printCmd.AddCommand(printPluginsCmd)

	printPluginsCmd.PersistentFlags().Bool("input", true, "print input plugins")
	printPluginsCmd.PersistentFlags().Bool("output", true, "print output plugins")
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints the info requested with argumets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var printMappingsCmd = &cobra.Command{
	Use:   "mappings",
	Short: "print mapping names to terminal",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := "mappings.yml"
		printMappingList(configFile)
	},
}

var printPluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "print available plugins in plugins path",
	Run: func(cmd *cobra.Command, args []string) {
		pluginsPath := "bin/plugins"
		printInputPlugins := cmd.Flags().Changed("input")
		printOutputPlugins := cmd.Flags().Changed("output")
		printAllPlugins := !(printInputPlugins || printOutputPlugins)

		if printAllPlugins || printInputPlugins {
			plugins, err := util.GetPluginNamesByType(pluginsPath, "input")
			if err != nil {
				panic("Failed to get input plugins")
			}
			fmt.Println("Input Plugins:")
			for i, plugin := range plugins {
				fmt.Printf("%4d %v\n", i+1, strings.TrimSuffix(plugin, ".so"))
			}
		}

		if printAllPlugins || printOutputPlugins {
			plugins, err := util.GetPluginNamesByType(pluginsPath, "output")
			if err != nil {
				panic("Failed to get input plugins")
			}
			fmt.Println("Output Plugins:")
			for i, plugin := range plugins {
				fmt.Printf("%4d %v\n", i+1, strings.TrimSuffix(plugin, ".so"))
			}
		}
	},
}

// PrintMappingList prints the mapping names in config file
func printMappingList(configFile string) {
	mappings := util.GetMappings(configFile)
	for i, mapping := range mappings {
		for name := range mapping {
			fmt.Printf("%6d %s\n", i, name)
		}
	}
	if len(mappings) == 0 {
		fmt.Printf("No mappings found\n")
	} else if len(mappings) == 1 {
		fmt.Printf("\nFound %d mapping\n", len(mappings))
	} else {
		fmt.Printf("\nFound %d mappings\n", len(mappings))
	}
}
