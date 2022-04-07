package server

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/conf"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/routes"
	"github.com/bluecolor/tractor/pkg/server/feedbackend"
	"github.com/bluecolor/tractor/pkg/tasks"
)

func startHttp(config conf.Config, exit chan error) {
	repository, err := repo.New(config.DB)
	if err != nil {
		exit <- err
		return
	}
	workerClient := tasks.NewClient(config.Worker)

	http.Handle("/", routes.BuildRoutes(repository, workerClient))
	exit <- http.ListenAndServe(":3000", nil)
}
func startFeedBackend(config conf.Config, exit chan error) {
	fb := feedbackend.New(config.FeedBackend)
	fb.Start(exit)
}

func Start(config conf.Config) error {
	exit := make(chan error)
	services := []func(conf.Config, chan error){
		startHttp,
		startFeedBackend,
	}

	for _, service := range services {
		go service(config, exit)
	}
	return <-exit
}
