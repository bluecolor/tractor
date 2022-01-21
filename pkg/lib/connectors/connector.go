package connectors

import "github.com/bluecolor/tractor/pkg/models"

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
