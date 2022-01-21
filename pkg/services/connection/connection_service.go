package connection

import (
	"encoding/json"
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
func (s *Service) OneConnection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	connection := models.Connection{}
	if err := s.repo.First(&connection, id).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
	}
	utils.RespondwithJSON(w, http.StatusOK, connection)
}
func (s *Service) FindConnections(w http.ResponseWriter, r *http.Request) {
	connections := []models.Connection{}
	result := s.repo.Find(&connections)
	if result.Error != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, result.Error)
	}
	utils.RespondwithJSON(w, http.StatusOK, connections)
}
func (s *Service) CreateConnection(w http.ResponseWriter, r *http.Request) {
	var connection models.Connection
	if err := json.NewDecoder(r.Body).Decode(&connection); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := s.repo.Create(&connection).Error; err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
}
func (s *Service) FindConnectionTypes(w http.ResponseWriter, r *http.Request) {
	connectionTypes := []models.ConnectionType{}
	result := s.repo.Find(&connectionTypes)
	if result.Error != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, result.Error)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, connectionTypes)
}
