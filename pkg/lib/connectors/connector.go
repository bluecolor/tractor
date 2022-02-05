package connectors

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
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
	FindDatasets(pattern string) ([]meta.Dataset, error)
}

type Input interface {
	Connector
	Read(p meta.ExtParams, w wire.Wire) error
}
type Output interface {
	Connector
	Write(e meta.ExtParams, w wire.Wire) error
}
