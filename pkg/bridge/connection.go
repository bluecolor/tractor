package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/models"
)

type Connection struct {
	model *models.Connection
}

func NewConnection(model *models.Connection) (c *Connection) {
	return &Connection{model: model}
}
func (c *Connection) Model() *models.Connection {
	return c.model
}
func (c *Connection) Connection() (*types.Connection, error) {
	config, err := GetConfig(c.model.Config)
	if err != nil {
		return nil, err
	}
	return &types.Connection{
		Name:           c.model.Name,
		ConnectionType: c.model.ConnectionType.Code,
		Config:         config,
		AsSource:       c.model.AsSource,
		AsTarget:       c.model.AsTarget,
	}, nil
}
