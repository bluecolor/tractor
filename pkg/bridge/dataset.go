package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/models"
)

type Dataset struct {
	model *models.Dataset
}

func NewDataset(model *models.Dataset) *Dataset {
	return &Dataset{
		model: model,
	}
}
func (d *Dataset) Model() *models.Dataset {
	return d.model
}
func (d *Dataset) WithFields(f []*models.Field) (*params.Dataset, error) {
	config, err := GetConfig(d.model.Config)
	if err != nil {
		return nil, err
	}
	fields, err := getFields(f)
	if err != nil {
		return nil, err
	}

	return &params.Dataset{
		Name:   d.model.Name,
		Config: config,
		Fields: fields,
	}, nil
}
