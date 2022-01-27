package dbutils

import (
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetParam(db *gorm.DB, name string) string {
	var param models.Param
	db.Where("name = ?", name).First(&param)
	return param.Value
}
func SetParam(db *gorm.DB, name string, value string) {
	param := models.Param{Name: name, Value: value}
	db.Save(&param)
}
func GetUniqueName(name string) string {
	return name + "_" + uuid.New().String()
}
