package connectors

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Connector interface {
	Connect() error
	Close() error
	GetDataset(options map[string]interface{}) (*types.Dataset, error)
	GetDatasets(options map[string]interface{}) ([]*types.Dataset, error)
	GetInfo(info string, options map[string]interface{}) (interface{}, error)
	Validate(config ConnectorConfig) error
}

type Base struct{}

func (b Base) Connect() error {
	return nil
}
func (b Base) Close() error {
	return nil
}

func (b Base) GetDataset(options map[string]interface{}) (*types.Dataset, error) {
	return nil, errors.New("not implemented")
}
func (b Base) GetDatasets(options map[string]interface{}) ([]*types.Dataset, error) {
	return nil, errors.New("not implemented")
}
func (b Base) GetInfo(info string, options map[string]interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}
func (b Base) Validate(config ConnectorConfig) error {
	return errors.New("not implemented")
}

type Input interface {
	Connector
	Read(d types.Dataset, w *wire.Wire) error
}
type Output interface {
	Connector
	Write(d types.Dataset, w *wire.Wire) error
}
