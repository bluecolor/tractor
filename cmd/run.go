package cmd

import (
	"fmt"
	"os"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run given mapping",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Mapping name is not given")
			os.Exit(1)
		}
		configFile := "mappings.yml"
		run(configFile, args[0])
	},
}

func run(configFile string, mapping string) {
	m, _ := util.GetMapping(configFile, mapping)

	if _, ok := m.Input["plugin"]; !ok {
		panic("Input plugin type is not given in mapping: " + mapping)
	}

	if _, ok := m.Output["plugin"]; !ok {
		panic("Output plugin type is not given in mapping: " + mapping)
	}
	inputPluginName := m.Input["plugin"].(string)
	outputPluginName := m.Output["plugin"].(string)

	inputPlugin, outputPlugin, err := util.GetPlugins("bin/plugins", inputPluginName, outputPluginName)

	if err != nil {
		panic(err)
	}

	inRunSymbol, err := inputPlugin.Lookup("Run")
	if err != nil {
		fmt.Printf("Failed find Run method in input plugin %s: %v\n", inputPluginName, err)
		os.Exit(1)
	}
	outRunSymbol, err := outputPlugin.Lookup("Run")
	if err != nil {
		fmt.Printf("Failed find Run method in output plugin %s: %v\n", outputPluginName, err)
		os.Exit(1)
	}
	inRunSymbol.(func(api.Config))(m.Input)
	outRunSymbol.(func(api.Config))(m.Output)
}
