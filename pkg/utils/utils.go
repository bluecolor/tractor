package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bluecolor/tractor/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RespondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func ErrorWithJSON(w http.ResponseWriter, code int, err error) {
	RespondwithJSON(w, code, map[string]string{"error": err.Error()})
}
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
