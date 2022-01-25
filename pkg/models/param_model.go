package models

type Param struct {
	Model
	Name  string `gorm:"size:200;not null;unique" json:"name"`
	Value string `gorm:"size:1000" json:"value"`
}
