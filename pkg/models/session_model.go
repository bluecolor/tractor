package models

import (
	"time"

	"gorm.io/datatypes"
)

type SessionStatus struct {
	Model
	Name string `json:"name"`
	Code string `json:"code"`
}

type Session struct {
	Model
	Extraction   *Extraction    `gorm:"foreignkey:ExtractionID" json:"extraction"`
	ExtractionID uint           `json:"extractionID"`
	Status       *SessionStatus `gorm:"foreignkey:StatusID" json:"status"`
	StatusID     uint           `json:"statusID"`
	StartedAt    *time.Time     `json:"startedAt"`
	FinishedAt   *time.Time     `json:"finishedAt"`
	WriteCount   int            `json:"writeCount"`
	Config       datatypes.JSON `gorm:"type:text" json:"config"`
}
