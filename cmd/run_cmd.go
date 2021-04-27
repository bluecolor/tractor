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
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

var progress bool

type runp struct {
	total   int
	read    int
	written int
	rpb     *mpb.Bar
	wpb     *mpb.Bar
}

func (rp *runp) init() {
	p := mpb.New(mpb.WithWidth(64))
	rp.rpb = p.AddBar(int64(rp.total),
		mpb.PrependDecorators(
			decor.Name("Read"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			),
		),
	)
	rp.wpb = p.AddBar(int64(rp.total),
		mpb.PrependDecorators(
			decor.Name("Write"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			),
		),
	)
}

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

	rp := &runp{}
	if progress {
		if counter, ok := inputPlugin.(tractor.Counter); ok {
			rp.total, err = counter.Count()
			if err != nil {
				println("Failed to get total count. Use without 'progress' flag")
				os.Exit(1)
			}
			rp.init()
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
	}(&wg)
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		outputPlugin.Write(wire)
		wg.Done()
		println("output done")
		wire.CloseFeed()
	}(&wg)
	wg.Add(1)

	go checkFeeds(wire, rp)

	wg.Wait()
	if progress {
		rp.rpb.Abort(true)
		rp.wpb.Abort(true)
	}

	duration := time.Since(start)
	fmt.Println("Duration:", duration)
}

func checkFeeds(wire tractor.Wire, rp *runp) {
	for f := range wire.ReadFeeds() {
		switch f.Type {
		case tractor.Progress:
			if progress {
				processProgressFeed(f, rp)
			}
		case tractor.Success:
			// println("Success", f.Sender)
		case tractor.Error:
			// println("Error", f.Sender)
		}
	}
}

func processProgressFeed(f tractor.Feed, rp *runp) {
	p := f.Content.(tractor.ProgressFeed)
	switch f.Sender {
	case tractor.InputPlugin:
		rp.read += p.Count()
		rp.rpb.IncrBy(p.Count())
	case tractor.OutputPlugin:
		rp.written += p.Count()
		rp.wpb.IncrBy(p.Count())
	}
}

func init() {
	runCmd.PersistentFlags().StringVar(&mapping, "mapping", "", "Mapping name")
	runCmd.PersistentFlags().BoolVar(&progress, "progress", false, "Show progress")
}
