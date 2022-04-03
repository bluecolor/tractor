package server

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/conf"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/routes"
	"github.com/bluecolor/tractor/pkg/tasks"
)

func Start(config conf.Config) error {
	repository, err := repo.New(config.DB)
	if err != nil {
		return err
	}
	workerClient := tasks.NewClient(config.Worker)

	http.Handle("/", routes.BuildRoutes(repository, workerClient))
	return http.ListenAndServe(":3000", nil)
}
