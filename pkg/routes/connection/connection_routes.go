package connection

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/connection"
	"github.com/go-chi/chi"
)

func BuildRoutes(repository *repo.Repository) *chi.Mux {
	service := connection.NewService(repository)
	router := chi.NewRouter()
	router.Mount("/", buildConnectionRoutes(service))
	router.Mount("/providers", buildProviderRoutes(service))
	router.Mount("/datasets", buildDatasetRoutes(service))
	return router
}

func buildConnectionRoutes(service *connection.Service) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", service.FindConnections)
	router.Get("/{id}", service.OneConnection)
	router.Get("/types", service.FindConnectionTypes)

	return router
}
func buildProviderRoutes(service *connection.Service) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", service.FindProviders)
	router.Get("/{id}", service.OneProvider)
	router.Delete("/{id}", service.DeleteProvider)
	router.Post("/", service.CreateProvider)
	return router
}
func buildDatasetRoutes(service *connection.Service) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", service.FindDatasets)
	router.Get("/{id}", service.OneDataset)
	router.Delete("/{id}", service.DeleteDataset)
	router.Post("/", service.CreateDataset)
	return router
}
