package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func MapToStruct(m map[string]interface{}, s interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s)
}
