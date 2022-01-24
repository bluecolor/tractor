package models

import "gorm.io/datatypes"

type ExtractionMode struct {
	Model
	Name string `json:"name"`
}

type Extraction struct {
	Model
	Name             string          `json:"name"`
	SourceDataset    datatypes.JSON  `gorm:"type:text" json:"sourceDataset"`
	TargetDataset    datatypes.JSON  `gorm:"type:text" json:"targetDataset"`
	ExtractionModeID uint            `json:"extractionModeID"`
	ExtractionMode   *ExtractionMode `gorm:"foreignkey:ExtractionModeID" json:"extractionMode"`
	FieldMappings    []FieldMapping  `json:"fieldMappings"`
}

type FieldMapping struct {
	Model
	SourceField  datatypes.JSON `gorm:"type:text" json:"sourceDataset"`
	TargetField  datatypes.JSON `gorm:"type:text" json:"targetDataset"`
	Extraction   *Extraction    `json:"extraction"`
	ExtractionID uint           `json:"extractionID"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}
