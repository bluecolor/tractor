package connection

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
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
	pattern := r.URL.Query().Get("pattern")
	connectionID := chi.URLParam(r, "connectionID")

	connection := models.Connection{}
	if err := s.repo.Preload("ConnectionType").First(&connection, connectionID).Error; err != nil {
		log.Error().Err(err).Msg("error getting connection")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	connectorCreator, ok := connectors.Connectors[connection.ConnectionType.Code]
	if !ok {
		err := errors.New("unsupported connection type " + connection.ConnectionType.Code)
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	connector, err := connectorCreator(connection.GetConfig())
	if err != nil {
		log.Error().Err(err).Msg("error creating connector")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	mfc, ok := connector.(connectors.MetaFetcher)
	if !ok {
		err := errors.New("connector does not implement metadata fetcher")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := mfc.Connect(); err != nil {
		log.Error().Err(err).Msg("error connecting to connector")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	datasets, err := mfc.FetchDatasets(pattern)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	if err := mfc.Close(); err != nil {
		log.Error().Err(err).Msg("error closing connector")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondwithJSON(w, http.StatusOK, datasets)
}
