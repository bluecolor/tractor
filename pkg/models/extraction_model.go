package models

import "gorm.io/datatypes"

type ExtractionMode struct {
	Model
	Name string `json:"name"`
	Code string `json:"code"`
}

type Extraction struct {
	Model
	Name             string          `json:"name"`
	SourceDatasetID  uint            `json:"sourceDatasetId"`
	SourceDataset    *Dataset        `gorm:"foreignkey:SourceDatasetID" json:"sourceDataset"`
	TargetDatasetID  uint            `json:"targetDatasetId"`
	TargetDataset    *Dataset        `gorm:"foreignkey:TargetDatasetID" json:"targetDataset"`
	ExtractionModeID uint            `json:"extractionModeID"`
	ExtractionMode   *ExtractionMode `gorm:"foreignkey:ExtractionModeID" json:"extractionMode"`
	FieldMappings    []FieldMapping  `json:"fieldMappings"`
}

type Dataset struct {
	Model
	Name         string         `json:"name"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
	ExtractionID uint           `json:"extractionID"`
	Extraction   *Extraction    `gorm:"foreignkey:ExtractionID" json:"extraction"`
}

type Field struct {
	Model
	Name       string         `json:"name"`
	Expression string         `json:"expression"`
	Type       string         `json:"type"`
	Config     datatypes.JSON `gorm:"type:text" json:"config"`
	DatasetID  uint           `json:"datasetID"`
	Dataset    *Dataset       `gorm:"foreignkey:DatasetID" json:"dataset"`
}

type FieldMapping struct {
	Model
	SourceField  *Field         `gorm:"type:text" json:"sourceDataset"`
	TargetField  *Field         `gorm:"type:text" json:"targetDataset"`
	Extraction   *Extraction    `json:"extraction"`
	ExtractionID uint           `json:"extractionID"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}
