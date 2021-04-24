package cmd

import (
	"fmt"
	"os"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/agent"
	cfg "github.com/bluecolor/tractor/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run given mapping",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	if mapping == "" {
		println("Mapping is not given")
		os.Exit(1)
	}
	conf := cfg.NewConfig()
	err := conf.LoadConfig(config)
	if err != nil {
		println("Failed to load config file")
		os.Exit(1)
	}
	m, err := conf.GetMapping(mapping)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	inputPlugin, err := validateAndGetInputPlugin(m.Input.Plugin, m.Input.Config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if initializer, ok := inputPlugin.(tractor.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}
	wire := agent.NewWire()
	err = inputPlugin.Read(wire)
	if err != nil {
		println(err.Error())
	}
}

func init() {
	runCmd.PersistentFlags().StringVar(&mapping, "mapping", "", "Mapping name")
}
