package routes

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/connection"
	"github.com/go-chi/chi"
)

func buildConnectionRoutes(repository *repo.Repository) *chi.Mux {
	service := connection.NewService(repository)
	router := chi.NewRouter()
	// connection routes
	router.Get("/", service.FindConnections)
	router.Get("/{id}", service.OneConnection)
	router.Get("/types", service.FindConnectionTypes)
	router.Post("/", service.CreateConnection)
	router.Post("/test", service.TestConnection)
	router.Delete("/{id}", service.DeleteConnection)
	router.Put("/{id}", service.UpdateConnection)
	router.Post("/{id}/dataset", service.GetDataset)
	router.Post("/info", service.GetInfo)

	// provider routes
	router.Get("/providers/types", service.FindProviderTypes)
	router.Get("/providers", service.FindProviders)
	router.Get("/providers/{id}", service.OneProvider)
	router.Delete("/providers/{id}", service.DeleteProvider)
	router.Post("/providers", service.CreateProvider)

	// dataset routes
	router.Get("/{connectionID}/datasets", service.FindDatasets)

	return router
}
