package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		Name:  "tractor",
		Usage: `🚜 tractor - data ingestion tool`,
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run",
				Action: func(c *cli.Context) error {
					println("Run")
					config := c.String("config")
					println(config)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "config file",
						Value:   "",
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
