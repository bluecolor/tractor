package models

import "gorm.io/datatypes"

type Extraction struct {
	Model
	Name               string         `json:"name"`
	SourceConnection   *Connection    `gorm:"foreignkey:SourceConnectionID" json:"sourceConnection"`
	SourceConnectionID uint           `json:"sourceConnectionId"`
	TargetConnection   *Connection    `gorm:"foreignkey:TargetConnectionID" json:"targetConnection"`
	TargetConnectionID uint           `json:"targetConnectionId"`
	SourceDatasetID    uint           `json:"sourceDatasetId"`
	SourceDataset      *Dataset       `gorm:"foreignkey:SourceDatasetID" json:"sourceDataset"`
	TargetDatasetID    uint           `json:"targetDatasetId"`
	TargetDataset      *Dataset       `gorm:"foreignkey:TargetDatasetID" json:"targetDataset"`
	FieldMappings      []FieldMapping `json:"fieldMappings"`
	Sessions           []*Session     `json:"sessions"`
}

type Dataset struct {
	Model
	Name         string         `json:"name"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
	ExtractionID uint           `json:"extractionID"`
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
	SourceField  *Field         `gorm:"type:text" json:"sourceField"`
	TargetField  *Field         `gorm:"type:text" json:"targetField"`
	Extraction   *Extraction    `json:"extraction"`
	ExtractionID uint           `json:"extractionID"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}

func (e *Extraction) GetSourceTargetFields() (s []*Field, t []*Field) {
	for _, fm := range e.FieldMappings {
		s = append(s, fm.SourceField)
		t = append(t, fm.TargetField)
	}
	return
}
