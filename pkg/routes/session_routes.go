package routes

import (
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/services/session"
	"github.com/go-chi/chi"
)

func buildSessionRoutes(repository *repo.Repository) *chi.Mux {
	service := session.NewService(repository)
	router := chi.NewRouter()
	router.Get("/", service.FindSessions)
	router.Delete("/{id}", service.DeleteSession)
	return router
}
