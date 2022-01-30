package connection

import (
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
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
	utils.RespondwithJSON(w, http.StatusOK, connection)
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
func (s *Service) TestConnection(w http.ResponseWriter, r *http.Request) {
	connection := models.Connection{}
	if err := json.NewDecoder(r.Body).Decode(&connection); err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	connectorConfig, err := connection.GetConnectorConfig()
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusBadRequest, err)
		return
	}
	connector, err := connectors.GetConnector(connection.ConnectionType.Code, connectorConfig)
	if err != nil {
		log.Error().Err(err).Msg("error getting connector")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := connector.Connect(); err != nil {
		log.Error().Err(err).Msg("error connecting")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := connector.Close(); err != nil {
		log.Error().Err(err).Msg("error closing")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, "success")
}
