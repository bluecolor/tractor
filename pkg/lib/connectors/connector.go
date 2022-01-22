package connectors

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/connectors/mysql"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/rs/zerolog/log"
)

type InputConnector interface {
}
type OutputConnector interface {
}

type Connector interface {
	InputConnector
	OutputConnector
	Connect() error
	Close() error
	FetchDatasets() ([]models.Dataset, error)
	FetchDatasetsWithPattern(pattern string) ([]models.Dataset, error)
}

func FetchDatasetsWithPattern(pattern string, connection models.Connection) ([]models.Dataset, error) {
	var connector Connector
	var err error
	log.Debug().Msgf("connection type code is %s", connection.ConnectionType.Code)
	switch connection.ConnectionType.Code {
	case "mysql":
		if connector, err = mysql.NewMySQLConnector(connection.GetConfig()); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported connection type " + connection.ConnectionType.Code)
	}
	if err := connector.Connect(); err != nil {
		return nil, err
	}
	datasets, err := connector.FetchDatasetsWithPattern(pattern)
	if err != nil {
		return nil, err
	}
	if err := connector.Close(); err != nil {
		return nil, err
	}
	return datasets, nil
}
