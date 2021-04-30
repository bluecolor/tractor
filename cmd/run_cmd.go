package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bluecolor/tractor/agent"
	"github.com/spf13/cobra"
)

var inputParams string
var outputParams string
var showProgress bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run given mapping",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {

	var session = &session{startTime: time.Now()}
	err := session.init()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	session.wire = agent.NewWire()

	wg, err := session.start()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	go session.listenFeeds()

	wg.Wait()
	session.end()
	fmt.Println("Duration:", session.duration())
}

func init() {
	runCmd.PersistentFlags().StringVar(&mapping, "mapping", "", "Mapping name")
	runCmd.PersistentFlags().BoolVar(&showProgress, "progress", false, "Show progress")
	runCmd.PersistentFlags().StringVar(&inputParams, "input-params", "", "Additional input plugin parameters")
	runCmd.PersistentFlags().StringVar(&outputParams, "output-params", "", "Additional output plugin parameters")
}
