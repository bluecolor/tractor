package models

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"gorm.io/datatypes"
)

type ConnectionType struct {
	Model
	Name          string         `gorm:"size:100;not null;unique" json:"name"`
	Code          string         `gorm:"size:100;not null;unique" json:"code"`
	ProviderTypes []ProviderType `gorm:"many2many:connection_type_provider_types"`
}
type Connection struct {
	Model
	Name             string          `gorm:"size:100;not null;unique" json:"name"`
	ConnectionTypeID uint            `json:"connectionTypeId"`
	ConnectionType   *ConnectionType `gorm:"foreignkey:ConnectionTypeID" json:"connectionType"`
	Config           datatypes.JSON  `gorm:"type:text" json:"config"`
	AsSource         bool            `gorm:"default:false" json:"asSource"`
	AsTarget         bool            `gorm:"default:false" json:"asTarget"`
}

func (c *Connection) GetConfig() connectors.ConnectorConfig {
	return connectors.ConnectorConfig(c.Config)
}

type ProviderType struct {
	Model
	Name            string           `gorm:"size:100;not null;unique" json:"name"`
	Code            string           `gorm:"size:100;not null;unique" json:"code"`
	ConnectionTypes []ConnectionType `gorm:"many2many:connection_type_provider_types"`
}
type Provider struct {
	Model
	Name           string         `gorm:"size:100;not null;unique" json:"name"`
	Config         datatypes.JSON `gorm:"type:text" json:"config"`
	ProviderTypeID uint           `json:"providerTypeId"`
	ProviderType   *ProviderType  `gorm:"foreignkey:ProviderTypeID" json:"providerType"`
}
type FileType struct {
	Model
	Name string `gorm:"size:100;not null;unique" json:"name"`
}