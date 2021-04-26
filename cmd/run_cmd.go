package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

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
	start := time.Now()

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
		err = initializer.Init(m.Input.Catalog)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}

	outputPlugin, err := validateAndGetOutputPlugin(m.Output.Plugin, m.Output.Config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if initializer, ok := outputPlugin.(tractor.Initializer); ok {
		var catalog *cfg.Catalog = m.Input.Catalog
		if catalog == nil {
			if disc, ok := inputPlugin.(tractor.Discoverer); ok {
				catalog, err = disc.Discover()
				if err != nil {
					println("Failed to discover catalog")
					println(err.Error())
				}
			}
		}
		err = initializer.Init(catalog)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}

	wire := agent.NewWire()

	var wg sync.WaitGroup
	go func(wg *sync.WaitGroup) {
		if err := inputPlugin.Read(wire); err != nil {
			println(err.Error())
		}
		wg.Done()
		wire.CloseData()
		println("input done")
	}(&wg)
	println("inputStarted")
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		outputPlugin.Write(wire)
		wg.Done()
		println("output done")
		wire.CloseFeed()
	}(&wg)
	println("outputStarted")
	wg.Add(1)

	go checkFeeds(wire)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("Duration:", duration)
}

func checkFeeds(wire tractor.Wire) {
	for f := range wire.ReadFeeds() {
		switch f.Type {
		case tractor.Progress:
			println("Progress ", f.Sender)
		case tractor.Success:
			println("Success ", f.Sender)
		case tractor.Error:
			println("Error ", f.Sender)
		}
	}
}

func init() {
	runCmd.PersistentFlags().StringVar(&mapping, "mapping", "", "Mapping name")
}
