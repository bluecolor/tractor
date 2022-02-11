package connectors

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Connector interface {
	Connect() error
	Close() error
}

type Base struct{}

func (c Base) Connect() error {
	return nil
}
func (c Base) Close() error {
	return nil
}

type MetaFinder interface {
	Connector
	FindDatasets(pattern string) ([]params.Dataset, error)
}

type Input interface {
	Connector
	Read(p params.SessionParams, w *wire.Wire) error
}
type Output interface {
	Connector
	Write(e params.SessionParams, w *wire.Wire) error
}
