package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/models"
)

type Field struct {
	model *models.Field
}

func NewField(model *models.Field) *Field {
	return &Field{
		model: model,
	}
}
func (f *Field) Model() *models.Field {
	return f.model
}
func (f *Field) Field() (*types.Field, error) {
	config, err := GetConfig(f.model.Config)
	if err != nil {
		return nil, err
	}

	return &types.Field{
		Name:   f.model.Name,
		Type:   types.FieldTypeFromString(f.model.Type),
		Config: config,
	}, nil
}

func getFields(fields []*models.Field) (output []*types.Field, err error) {
	output = make([]*types.Field, len(fields))
	for i, f := range fields {
		output[i], err = NewField(f).Field()
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func getFieldMappings(m []models.FieldMapping) (output []types.FieldMapping, err error) {
	output = make([]types.FieldMapping, len(m))
	for i, fm := range m {
		source, err := NewField(fm.SourceField).Field()
		if err != nil {
			return nil, err
		}
		target, err := NewField(fm.TargetField).Field()
		if err != nil {
			return nil, err
		}
		config, err := GetConfig(fm.Config)
		if err != nil {
			return nil, err
		}
		output[i] = types.FieldMapping{
			SourceField: source,
			TargetField: target,
			Config:      config,
		}
	}
	return output, nil
}
