package models

type ExtractionMode struct {
	Model
	Name string `json:"name"`
}

type Extraction struct {
	Model
	Name             string          `json:"name"`
	SourceDatasetID  uint            `json:"sourceDatasetID"`
	SourceDataset    *Dataset        `gorm:"foreignkey:SourceDatasetID" json:"sourceDataset"`
	TargetDatasetID  uint            `json:"targetDatasetID"`
	TargetDataset    *Dataset        `gorm:"foreignkey:TargetDatasetID" json:"targetDataset"`
	ExtractionModeID uint            `json:"extractionModeID"`
	ExtractionMode   *ExtractionMode `gorm:"foreignkey:ExtractionModeID" json:"extractionMode"`
}

type Mapping struct {
	Model
	TargetFieldID uint   `json:"targetFieldID"`
	TargetField   *Field `gorm:"foreignkey:TargetFieldID" json:"targetField"`
	SourceFieldID uint   `json:"sourceFieldID"`
	SourceField   *Field `gorm:"foreignkey:SourceFieldID" json:"sourceField"`
	Expression    string `gorm:"default:NULL" json:"expression"` // reserved for future use
}
