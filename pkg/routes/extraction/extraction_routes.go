package extraction

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/extraction"
	"github.com/go-chi/chi"
)

func BuildRoutes(repository *repo.Repository) *chi.Mux {
	service := extraction.NewService(repository)
	router := chi.NewRouter()
	router.Get("/", service.FindExtractions)
	router.Get("/{id}", service.OneExtraction)
	router.Delete("/{id}", service.DeleteExtraction)
	router.Post("/", service.CreateExtraction)
	return router
}
