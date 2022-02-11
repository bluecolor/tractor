package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/models"
)

type Extraction struct {
	model models.Extraction
}

func NewExtraction(model models.Extraction) (e *Extraction) {
	return &Extraction{model: model}
}

func (e *Extraction) Model() models.Extraction {
	return e.model
}

func (e *Extraction) ExtParams() (p params.ExtParams, err error) {
	sourceFields, targetFields := e.model.GetSourceTargetFields()
	inputDataset, err := NewDataset(e.model.SourceDataset).DatasetWithFields(sourceFields)
	if err != nil {
		return p, err
	}
	outputDataset, err := NewDataset(e.model.SourceDataset).DatasetWithFields(targetFields)
	if err != nil {
		return p, err
	}
	fm, err := getFieldMappings(e.model.FieldMappings)
	if err != nil {
		return p, err
	}
	p = params.ExtParams{}.
		WithInputDataset(inputDataset).
		WithOutputDataset(outputDataset).
		WithFieldMappings(fm)
	return
}
