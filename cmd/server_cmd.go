package cmd

import (
	"github.com/bluecolor/tractor/pkg/server"
	"github.com/urfave/cli/v2"
)

func runServerStartCmd(c *cli.Context) error {
	return server.Start(config)
}
