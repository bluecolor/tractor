package connection

import (
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
)

func (s *Service) OneDataset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dataset := models.Dataset{}
	if err := s.repo.First(&dataset, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	utils.RespondwithJSON(w, http.StatusOK, dataset)
}
func (s *Service) FindDatasets(w http.ResponseWriter, r *http.Request) {
	datasets := []models.Dataset{}
	if err := s.repo.Find(&datasets).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, datasets)
}
func (s *Service) CreateDataset(w http.ResponseWriter, r *http.Request) {
	var dataset models.Dataset
	if err := json.NewDecoder(r.Body).Decode(&dataset); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := s.repo.Create(&dataset).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, dataset)
}
func (s *Service) DeleteDataset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dataset := models.Dataset{}
	if err := s.repo.First(&dataset, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.Delete(&dataset).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, dataset)
}
func (s *Service) FetchDatasets(w http.ResponseWriter, r *http.Request) {
	datasets := []models.Dataset{}
	utils.RespondwithJSON(w, http.StatusOK, datasets)
}
