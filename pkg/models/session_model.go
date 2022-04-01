package models

import (
	"time"

	"gorm.io/datatypes"
)

type Session struct {
	Model
	Extraction   *Extraction    `gorm:"foreignkey:ExtractionID" json:"extraction"`
	ExtractionID uint           `json:"extractionID"`
	Status       string         `json:"status"`
	StartedAt    *time.Time     `json:"startedAt"`
	FinishedAt   *time.Time     `json:"finishedAt"`
	WriteCount   int            `json:"writeCount"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}
