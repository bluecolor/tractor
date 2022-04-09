package server

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/routes"
	"github.com/bluecolor/tractor/pkg/tasks"
)

func Start(config conf.Config) error {
	repository, err := repo.New(config.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create repo")
	}
	workerClient := tasks.NewClient(config.Worker)

	http.Handle("/", routes.BuildRoutes(repository, workerClient))
	return http.ListenAndServe(":3000", nil)
}
