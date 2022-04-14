package models

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"gorm.io/datatypes"
)

type ConnectionType struct {
	Model
	Name string `gorm:"size:100;not null;unique" json:"name"`
	Code string `gorm:"size:100;not null;unique" json:"code"`
}
type Connection struct {
	Model
	Name             string          `gorm:"size:100;not null;unique" json:"name"`
	ConnectionTypeID uint            `json:"connectionTypeId"`
	ConnectionType   *ConnectionType `fake:"skip" gorm:"foreignkey:ConnectionTypeID" json:"connectionType"`
	Config           datatypes.JSON  `fake:"skip" gorm:"type:text" json:"config"`
	AsSource         bool            `gorm:"default:true" json:"asSource"`
	AsTarget         bool            `gorm:"default:true" json:"asTarget"`
}

func (c *Connection) GetConnectorConfig() (connectors.ConnectorConfig, error) {
	configMap := connectors.ConnectorConfig{}
	if err := json.Unmarshal(c.Config, &configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

type Provider struct {
	Model
	Name string `gorm:"size:100;not null;unique" json:"name"`
	Code string `gorm:"size:100;not null;unique" json:"name"`
}
type FileType struct {
	Model
	Name string `gorm:"size:100;not null;unique" json:"name"`
	Code string `gorm:"size:100;not null;unique" json:"code"`
}
