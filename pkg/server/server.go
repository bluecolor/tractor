package server

import (
	"net/http"

	"github.com/go-redis/redis/v8"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Worker.FeedProcessor.CacheAddr,
	})

	http.Handle("/", routes.BuildRoutes(repository, workerClient, redisClient))
	return http.ListenAndServe(":3000", nil)
}
