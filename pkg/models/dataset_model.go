package models

import "gorm.io/datatypes"

type Dataset struct {
	Model
	Name         string         `json:"name"`
	Fields       []Field        `gorm:"foreignkey:DatasetID" json:"fields"`
	ConnectionID uint           `json:"connectionID"`
	Connection   *Connection    `gorm:"foreignkey:ConnectionID" json:"connection"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}

type Field struct {
	Model
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	DatasetID uint           `json:"datasetID"`
	Dataset   *Dataset       `gorm:"foreignkey:DatasetID" json:"dataset"`
	Config    datatypes.JSON `gorm:"type:text" json:"config"`
}
