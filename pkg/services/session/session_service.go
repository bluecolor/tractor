package session

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/utils"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) FindSessions(w http.ResponseWriter, r *http.Request) {
	sessions := []models.Session{}
	result := s.repo.
		Preload("Extraction").
		Preload("Extraction.SourceDataset").
		Preload("Extraction.SourceDataset.Connection").
		Preload("Extraction.TargetDataset").
		Preload("Extraction.TargetDataset.Connection").
		Find(&sessions)
	if result.Error != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, result.Error)
	}
	utils.RespondwithJSON(w, http.StatusOK, sessions)
}
