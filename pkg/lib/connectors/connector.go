package connectors

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Connector interface {
	Connect() error
	Close() error
}

type BaseConnector struct{}

func (c BaseConnector) Connect() error {
	return nil
}
func (c BaseConnector) Close() error {
	return nil
}

type MetaFinder interface {
	Connector
	FindDatasets(pattern string) ([]meta.Dataset, error)
}

type InputConnector interface {
	Connector
	Read(p meta.ExtParams, w wire.Wire) error
}
type OutputConnector interface {
	Connector
	Write(e meta.ExtParams, w wire.Wire) error
}
