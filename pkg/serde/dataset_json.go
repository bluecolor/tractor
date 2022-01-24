package serde

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/models"
)

type FetchDatasetOutput models.Dataset

func (d *FetchDatasetOutput) MarshalJSON() ([]byte, error) {
	type field struct {
		Name   string          `json:"name"`
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}
	type dataset struct {
		Fields []field         `json:"fields"`
		Config json.RawMessage `json:"config"`
	}
	ds := dataset{
		Fields: make([]field, len(d.Fields)),
	}
	for _, f := range d.Fields {
		ds.Fields = append(ds.Fields, field{
			Name:   f.Name,
			Type:   f.Type,
			Config: json.RawMessage(f.Config),
		})
	}

	return json.Marshal(&struct {
		Name         string             `json:"name"`
		Fields       []field            `json:"fields"`
		ConnectionID uint               `json:"connectionID"`
		Connection   *models.Connection `json:"connection"`
		Config       json.RawMessage    `json:"config"`
	}{})
}
