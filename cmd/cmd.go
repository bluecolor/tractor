package cmd

import (
	"os"

	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/session"
	"github.com/urfave/cli/v2"
)

func execRun(c *cli.Context) (err error) {
	configPath := c.String("config")
	progress := c.Bool("progress")
	params := c.String("params")
	mapping := c.String("mapping")
	if configPath == "" {
		return cli.Exit("config is required", 1)
	}
	if mapping == "" {
		return cli.Exit("mapping is required", 1)
	}
	conf := config.NewConfig()
	err = conf.Load(configPath)
	if err != nil {
		return
	}
	s, err := session.NewSession(conf, mapping, progress, params)
	if err != nil {
		return
	}
	err = s.Run()
	return
}

func Run() {
	app := &cli.App{
		Name:  "tractor",
		Usage: `🚜 tractor - data ingestion tool`,
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run",
				Action: execRun,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "config file",
						Value:   "",
					},
					&cli.BoolFlag{
						Name:  "progress",
						Usage: "show progress",
						Value: false,
					},
					&cli.StringFlag{
						Name:    "params",
						Aliases: []string{"p"},
						Usage:   "additional parameters",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "mapping",
						Aliases: []string{"m"},
						Usage:   "mapping name",
						Value:   "",
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
