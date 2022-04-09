package routes

import (
	"github.com/bluecolor/tractor/pkg/services/ws"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v7"
)

func buildWSRoutes(client *redis.Client) *chi.Mux {
	service := ws.NewService(client)
	router := chi.NewRouter()
	router.Get("/session/feeds", service.SubSessionFeeds)
	return router
}
