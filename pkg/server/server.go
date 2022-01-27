package server

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/conf"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/routes"
)

func Start(config conf.Config) error {
	repository, err := repo.NewRepository(config.DB)
	if err != nil {
		return err
	}
	http.Handle("/", routes.BuildRoutes(repository))
	return http.ListenAndServe(":3000", nil)
}
