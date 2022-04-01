package models

import (
	"time"
)

var Models = []interface{}{}

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func init() {
	Models = append(Models,
		&Connection{},
		&ConnectionType{},
		&Extraction{},
		&Dataset{},
		&Field{},
		&FileType{},
		&Provider{},
		&ProviderType{},
		&Param{},
		&Session{},
	)
}
