package commands

import (
	"fmt"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/logging"
	"github.com/bluecolor/tractor/util"
	c "github.com/bluecolor/tractor/util/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run given mapping",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		logging.Fatal("Mapping name is not given")
	}
	file := viper.GetString(c.TractorMappingsFile)
	mapping, err := util.GetMapping(file, args[0], true)
	if err != nil {
		logging.Fatal("Can not read mapping", err)
	}
	ip, op, err := util.GetMappingPlugins(mapping)
	if err != nil {
		logging.Fatal("Failed to get plugins", err)
	}
	iconf, oconf, err := util.GetConfigs(mapping)
	if err != nil {
		logging.Fatal("Failed to get configs from mapping")
	}

	wire := api.Wire{
		Feed:     make(chan *api.Feed, viper.GetInt(c.TractorChannelBufferSize)),
		Metadata: make(chan *api.Metadata, viper.GetInt(c.TractorChannelBufferSize)),
		Data:     make(chan *api.Data, viper.GetInt(c.TractorChannelBufferSize)),
	}

	var wg sync.WaitGroup
	go ip.Run(&wg, iconf, &wire)
	wg.Add(1)
	go op.Run(&wg, oconf, &wire)
	wg.Add(1)

	go func(ch chan *api.Feed) {
		var readCount, writeCount = 0, 0
		for feed := range ch {
			if feed.Type == api.ReadCountFeed || feed.Type == api.WriteCountFeed {
				switch feed.Type {
				case api.ReadCountFeed:
					readCount += feed.Content.(int)
				case api.WriteCountFeed:
					writeCount += feed.Content.(int)
				}
				fmt.Printf("\033[2K\rReadCount: %d WriteCount: %d", readCount, writeCount)
			}
		}
	}(wire.Feed)

	wg.Wait()
	close(wire.Feed)
	fmt.Println("\nDone")
}
