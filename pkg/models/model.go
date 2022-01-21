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
		&Dataset{},
		&Extraction{},
		&Field{},
		&FileType{},
		&Mapping{},
		&Provider{},
		&ProviderType{},
		&Param{},
	)
}
