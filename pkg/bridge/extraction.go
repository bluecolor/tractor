package bridge

import (
	"github.com/bluecolor/tractor/pkg/models"
)

type Extraction struct {
	model *models.Extraction
}

func NewExtraction(model *models.Extraction) (e *Extraction) {
	return &Extraction{model: model}
}

func (e *Extraction) Model() *models.Extraction {
	return e.model
}
