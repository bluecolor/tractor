package cmd

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func runConnectRepo(c *cli.Context) (err error) {
	log.Info().Msg("connecting to the repository")
	repository, err = repo.NewRepository(config.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the repository")
	}
	log.Info().Msg("connected to the repository")
	return
}
func runDbMigrateCmd(c *cli.Context) (err error) {
	log.Info().Msg("migrating database")
	err = repository.Migrate()
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate database")
		return
	}
	log.Info().Msg("migrated database")
	return nil
}
func runDbDropCmd(c *cli.Context) (err error) {
	log.Info().Msg("dropping tables")
	err = repository.Drop()
	if err != nil {
		log.Error().Err(err).Msg("failed to drop tables")
		return
	}
	log.Info().Msg("dropped tables")
	return
}
func runSeedCmd(c *cli.Context) (err error) {
	log.Info().Msg("seeding database")
	reset := c.Bool("reset")
	err = repository.Seed(config.App.SeedPath, reset)
	if err != nil {
		log.Error().Err(err).Msg("failed to drop tables")
		return
	}
	log.Info().Msg("seeded database")
	return
}
func runResetCmd(c *cli.Context) (err error) {
	log.Info().Msg("resetting database")
	err = repository.Seed(config.App.SeedPath, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to reset database")
		return
	}
	log.Info().Msg("database is reset")
	return
}
