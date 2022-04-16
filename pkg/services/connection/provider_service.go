package connection

import (
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
)

func (s *Service) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var provider models.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := s.repo.Create(&provider).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, provider)
}
func (s *Service) FindProviders(w http.ResponseWriter, r *http.Request) {
	providers := []models.Provider{}
	if err := s.repo.Find(&providers).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, providers)
}
func (s *Service) OneProvider(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var provider models.Provider = models.Provider{}
	if err := s.repo.First(&provider, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, provider)
}
func (s *Service) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	provider := models.Provider{}
	if err := s.repo.First(&provider, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.repo.Delete(&provider).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, provider)
}
