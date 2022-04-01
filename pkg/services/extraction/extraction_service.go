package extraction

import (
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/tasks"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
)

type Service struct {
	repo   *repo.Repository
	client *tasks.Client
}

func NewService(repo *repo.Repository, client *tasks.Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
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
	exts := []models.Extraction{}
	result := s.repo.
		Preload("SourceDataset").
		Preload("SourceDataset.Connection").
		Preload("TargetDataset.Connection").
		Preload("TargetDataset").
		Find(&exts)
	if result.Error != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, result.Error)
	}
	utils.RespondwithJSON(w, http.StatusOK, exts)
}
func (s *Service) DeleteExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ext := models.Extraction{}
	if err := s.repo.First(&ext, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.Delete(&ext).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, ext)
}
func (s *Service) CreateExtraction(w http.ResponseWriter, r *http.Request) {
	extraction := models.Extraction{}
	if err := json.NewDecoder(r.Body).Decode(&extraction); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := s.repo.Create(&extraction).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, extraction)
}
func (s *Service) RunExtraction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	extraction := models.Extraction{}
	if err := s.repo.First(&extraction, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	session := extraction.NewSession()
	if err := s.repo.Create(session).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	ses, err := tasks.GetSession(session)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	task, err := tasks.NewSessionRunTask(ses)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if _, err := s.client.Enqueue(task); err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, session)
}
