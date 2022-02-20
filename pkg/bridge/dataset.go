package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
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
func (d *Dataset) WithFields(f []*models.Field) (*types.Dataset, error) {
	config, err := GetConfig(d.model.Config)
	if err != nil {
		return nil, err
	}
	fields, err := getFields(f)
	if err != nil {
		return nil, err
	}

	return &types.Dataset{
		Name:   d.model.Name,
		Config: config,
		Fields: fields,
	}, nil
}
