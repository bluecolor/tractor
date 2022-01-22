package connection

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/connection"
	"github.com/go-chi/chi"
)

func BuildRoutes(repository *repo.Repository) *chi.Mux {
	service := connection.NewService(repository)
	router := chi.NewRouter()
	// connection routes
	router.Get("/", service.FindConnections)
	router.Get("/{id}", service.OneConnection)
	router.Get("/types", service.FindConnectionTypes)
	router.Post("/", service.CreateConnection)

	// provider routes
	router.Get("/providers", service.FindProviders)
	router.Get("/providers/{id}", service.OneProvider)
	router.Delete("/providers/{id}", service.DeleteProvider)
	router.Post("/providers", service.CreateProvider)

	// dataset routes
	router.Get("/{connectionID}/datasets/fetch", service.FetchDatasets)

	// router.Get("/{connectionID}/datasets", service.FindDatasets)
	// router.Get("/{connectionID}/datasets/{id}", service.OneDataset)
	// router.Delete("/{connectionID}/datasets/{id}", service.DeleteDataset)
	// router.Post("/{connectionID}/datasets", service.CreateDataset)

	return router
}
