package cmd

import (
	"fmt"
	"os"
	"plugin"
	"sync"

	"github.com/bluecolor/tractor/api/message"
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

func getRunMethod(plug *plugin.Plugin) (func(*sync.WaitGroup, []byte, chan message.Message), error) {
	symbol, err := plug.Lookup("Run")
	if err != nil {
		return nil, err
	}
	return symbol.(func(*sync.WaitGroup, []byte, chan message.Message)), nil
}

func run(configFile string, mapping string) {
	m, _ := util.GetMapping(configFile, mapping)

	if _, ok := m.Input["plugin"]; !ok {
		panic("Input plugin type is not given in mapping: " + mapping)
	}

	if _, ok := m.Output["plugin"]; !ok {
		panic("Output plugin type is not given in mapping: " + mapping)
	}
	iplugn := m.Input["plugin"].(string)
	oplugn := m.Output["plugin"].(string)

	iplug, oplug, err := util.GetPlugins("bin/plugins", iplugn, oplugn)
	if err != nil {
		panic(err)
	}

	irun, err := getRunMethod(iplug)
	if err != nil {
		fmt.Printf("Failed find Run method in input plugin %s: %v\n", iplugn, err)
		os.Exit(1)
	}
	orun, err := getRunMethod(oplug)
	if err != nil {
		fmt.Printf("Failed find Run method in output plugin %s: %v\n", oplugn, err)
		os.Exit(1)
	}

	iconf, oconf, err := getIOConf(m)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	channel := make(chan message.Message, 1000) // todo buffer from .env

	go irun(&wg, iconf, channel)
	wg.Add(1)
	go orun(&wg, oconf, channel)
	wg.Add(1)

	wg.Wait()
}
