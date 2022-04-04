package extraction

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/extraction"
	"github.com/bluecolor/tractor/pkg/tasks"
	"github.com/go-chi/chi"
)

func BuildRoutes(repository *repo.Repository, workerClient *tasks.Client) *chi.Mux {
	service := extraction.NewService(repository, workerClient)
	router := chi.NewRouter()
	router.Get("/", service.FindExtractions)
	router.Get("/{id}", service.OneExtraction)
	router.Delete("/{id}", service.DeleteExtraction)
	router.Post("/", service.CreateExtraction)
	router.Post("/{id}/run", service.RunExtraction)
	return router
}
