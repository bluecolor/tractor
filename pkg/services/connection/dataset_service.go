package connection

import (
	"errors"
	"net/http"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func (s *Service) FindDatasets(w http.ResponseWriter, r *http.Request) {
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
	connectorConfig, err := connection.GetConnectorConfig()
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	connector, err := connectorCreator(connectorConfig)
	if err != nil {
		log.Error().Err(err).Msg("error creating connector")
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	mfc, ok := connector.(connectors.MetaFinder)
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
	datasets, err := mfc.FindDatasets(pattern)
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
