package connectors

import (
	"github.com/bluecolor/tractor/lib/wire"
	"github.com/bluecolor/tractor/pkg/lib/cat/meta"
)

type Connector interface {
	Connect() error
	Close() error
}

type MetaFetcher interface {
	Connector
	FetchDatasets(pattern string) ([]meta.Dataset, error)
}

type InputConnector interface {
	Connector
	Read(w *wire.Wire) error
}
type OutputConnector interface {
	Connector
	Write(w *wire.Wire) error
}
