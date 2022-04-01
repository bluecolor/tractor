package models

import "gorm.io/datatypes"

type Extraction struct {
	Model
	Name            string         `json:"name"`
	SourceDatasetID uint           `json:"sourceDatasetId"`
	SourceDataset   *Dataset       `gorm:"foreignkey:SourceDatasetID" json:"sourceDataset"`
	TargetDatasetID uint           `json:"targetDatasetId"`
	TargetDataset   *Dataset       `gorm:"foreignkey:TargetDatasetID" json:"targetDataset"`
	Sessions        []*Session     `json:"sessions"`
	Config          datatypes.JSON `gorm:"type:text" json:"config"`
}

type Dataset struct {
	Model
	Name         string         `json:"name"`
	Connection   *Connection    `gorm:"foreignkey:ConnectionID" json:"connection"`
	ConnectionID uint           `json:"connectionId"`
	ExtractionID uint           `json:"extractionId"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
	Fields       []*Field       `json:"fields"`
}

type Field struct {
	Model
	Order      int            `json:"order"`
	Name       string         `json:"name"`
	Expression string         `json:"expression"`
	Type       string         `json:"type"`
	Config     datatypes.JSON `gorm:"type:text" json:"config"`
	DatasetID  uint           `json:"datasetId"`
	Dataset    *Dataset       `gorm:"foreignkey:DatasetID" json:"dataset"`
}

func (e *Extraction) NewSession() *Session {
	return &Session{
		ExtractionID: e.ID,
		Status:       "pending",
	}
}
