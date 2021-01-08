package cmd

import (
	"fmt"
	"os"

	"github.com/bluecolor/tractor/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

func getIOConf(m *util.Mapping) ([]byte, []byte, error) {
	iconf, err := yaml.Marshal(m.Input)
	if err != nil {
		return nil, nil, err
	}
	oconf, err := yaml.Marshal(m.Output)
	if err != nil {
		return nil, nil, err
	}
	return iconf, oconf, nil
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

	iconf, oconf, err := getIOConf(m)
	if err != nil {
		panic(err)
	}

	inRunSymbol.(func([]byte))(iconf)
	outRunSymbol.(func([]byte))(oconf)
}
