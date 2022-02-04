package extraction

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) OneExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ext := models.Extraction{}
	if err := s.repo.First(&ext, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	utils.RespondwithJSON(w, http.StatusOK, ext)
}

func (s *Service) FindExtractions(w http.ResponseWriter, r *http.Request) {
}
