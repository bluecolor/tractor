package session

import (
	"net/http"
	"strings"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/morkid/paginate"
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
	model := s.repo.
		Joins("Extraction").
		Preload("Extraction.SourceDataset").
		Preload("Extraction.SourceDataset.Connection").
		Preload("Extraction.TargetDataset").
		Preload("Extraction.TargetDataset.Connection").
		Model(&[]models.Session{})

	if r.URL.Query().Get("statuses") != "" {
		statuses := strings.Split(r.URL.Query().Get("statuses"), ",")
		model = model.Where("Extraction.status IN (?)", statuses)
	}
	if r.URL.Query().Get("extraction") != "" {
		model = model.Where("Extraction.id = ?", r.URL.Query().Get("extraction"))
	}
	if r.URL.Query().Get("q") != "" {
		model = model.Where("Extraction.name LIKE ?", "%"+r.URL.Query().Get("q")+"%")
	}

	result := paginate.New().Response(model, r, &[]models.Session{})
	utils.RespondwithJSON(w, http.StatusOK, result)
}

func (s *Service) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	session := models.Session{}
	if err := s.repo.First(&session, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.Delete(&session).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, session)
}
