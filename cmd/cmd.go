package cmd

import (
	"os"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	config     conf.Config
	repository *repo.Repository
)

func init() {
	var err error
	config, err = conf.LoadConfig()
	if err != nil {
		panic(err)
	}
	setupLogger(config.Log)
}
func setupLogger(conf conf.Log) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	return nil
}

func Run() {
	app := &cli.App{
		Name:  "tractor",
		Usage: `🚜 tractor - data ingestion tool`,
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "server command",
				Subcommands: []*cli.Command{
					{
						Name:   "start",
						Usage:  "start tractor server",
						Action: runServerStartCmd,
					},
				},
			},
			{
				Name:  "worker",
				Usage: "worker command",
				Subcommands: []*cli.Command{
					{
						Name:   "start",
						Usage:  "start tractor worker",
						Action: runWorkerStartCmd,
					},
				},
			},
			{
				Name:   "db",
				Usage:  "db command",
				Before: runConnectRepo,
				Subcommands: []*cli.Command{
					{
						Name:   "migrate",
						Usage:  "migrate",
						Action: runDbMigrateCmd,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:        "reset",
								Aliases:     []string{"r"},
								Usage:       "reset database",
								DefaultText: "false",
							},
						},
					},
					{
						Name:   "drop",
						Usage:  "drop",
						Action: runDbDropCmd,
					},
					{
						Name:   "seed",
						Usage:  "seed",
						Action: runSeedCmd,
					},
					{
						Name:   "reset",
						Usage:  "reset",
						Action: runResetCmd,
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
