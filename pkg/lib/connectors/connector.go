package connectors

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
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
	FindDatasets(options map[string]interface{}) ([]types.Dataset, error)
}
type FieldFinder interface {
	Connector
	FindFields(options map[string]interface{}) ([]types.Field, error)
}

type RequestResolver interface {
	GetResolvers() map[string]func(map[string]interface{}) (interface{}, error)
	Resolve(request string, body map[string]interface{}) (interface{}, error)
}

type Input interface {
	Connector
	Read(p types.SessionParams, w *wire.Wire) error
}
type Output interface {
	Connector
	Write(e types.SessionParams, w *wire.Wire) error
}
