package cmd

import (
	"fmt"
	"log"
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
			log.Fatalln("Mapping name is not given")
		}
		run(util.MappingsFile, args[0])
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

func run(mappingsFile string, mapping string) {
	m, _ := util.GetMapping(mappingsFile, mapping)

	if _, ok := m.Input["plugin"]; !ok {
		log.Fatalln("Input plugin type is not given in mapping: " + mapping)
	}

	if _, ok := m.Output["plugin"]; !ok {
		log.Fatalln("Output plugin type is not given in mapping: " + mapping)
	}

	iplugin := m.Input["plugin"].(string)
	oplugin := m.Output["plugin"].(string)

	iplug, oplug, err := util.GetPlugins("bin/plugins", iplugin, oplugin)
	if err != nil {
		panic(err)
	}

	irun, err := getRunMethod(iplug)
	if err != nil {
		fmt.Printf("Failed find Run method in input plugin %s: %v\n", iplugin, err)
		os.Exit(1)
	}
	orun, err := getRunMethod(oplug)
	if err != nil {
		fmt.Printf("Failed find Run method in output plugin %s: %v\n", oplugin, err)
		os.Exit(1)
	}

	iconf, oconf, err := getIOConf(m)
	if err != nil {
		panic(err)
	}

	channel := make(chan message.Message, util.GetChannelBuffer())

	var wg sync.WaitGroup
	go irun(&wg, iconf, channel)
	wg.Add(1)
	go orun(&wg, oconf, channel)
	wg.Add(1)

	wg.Wait()
}
