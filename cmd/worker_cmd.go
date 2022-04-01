package cmd

import (
	"github.com/bluecolor/tractor/pkg/tasks"
	"github.com/urfave/cli/v2"
)

func runWorkerStartCmd(c *cli.Context) error {
	return tasks.NewWorker(config.Worker).Start()
}
