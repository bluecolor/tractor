package meta

import (
	"fmt"
	"time"
)

const (
	TimeoutKey        = "timeout"
	InputDatasetKey   = "input_dataset"
	OutputDatasetKey  = "output_dataset"
	ParallelKey       = "parallel"
	ExtractionModeKey = "extraction_mode"
	FieldMappingsKey  = "field_mappings"
	BufferSizeKey     = "buffer_size"
	DefaultTimeOut    = time.Second * 60 * 10 // todo from env 10 minutes
)

type ExtParams map[string]interface{}

func (p ExtParams) WithTimeout(timeout time.Duration) ExtParams {
	p[TimeoutKey] = timeout
	return p
}
func (p ExtParams) WithInputDataset(dataset Dataset) ExtParams {
	p[InputDatasetKey] = dataset
	return p
}
func (p ExtParams) WithInputParallel(parallel int) ExtParams {
	p[InputDatasetKey].(Dataset).Config[ParallelKey] = parallel
	return p
}
func (p ExtParams) WithOutputDataset(dataset Dataset) ExtParams {
	p[OutputDatasetKey] = dataset
	return p
}
func (p ExtParams) WithExtractionModeString(mode string) ExtParams {
	p[ExtractionModeKey] = ExtractionModeFromString(mode)
	return p
}
func (p ExtParams) WithFieldMappings(mappings []FieldMapping) ExtParams {
	p[FieldMappingsKey] = mappings
	return p
}
func (p ExtParams) GetTimeout() time.Duration {
	if timeout, ok := p[TimeoutKey]; ok {
		if t, ok := timeout.(time.Duration); ok {
			return t
		}
	}
	return DefaultTimeOut
}
func (p ExtParams) GetExtractionMode() ExtractionMode {
	if mode, ok := p[ExtractionModeKey]; ok {
		if m, ok := mode.(ExtractionMode); ok {
			return m
		}
	}
	return ExtractionModeCreate
}
func (p ExtParams) GetInputDataset() *Dataset {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return &d
		}
	}
	return nil
}
func (p ExtParams) GetOutputDataset() *Dataset {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return &d
		}
	}
	return nil
}
func (p ExtParams) GetInputParallel() int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(ParallelKey, 1)
		}
	}
	return 1
}
func (p ExtParams) GetOutputParallel() int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(ParallelKey, 1)
		}
	}
	return 1
}
func (p ExtParams) GetInputBufferSize() int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, 1000) // todo default from env
		}
	}
	return 1000
}
func (p ExtParams) GetOutputBufferSize() int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, 1000) // todo default from env
		}
	}
	return 1000
}
func (p ExtParams) GetInputBufferSizeOrThis(this int) int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, this)
		}
	}
	return this
}
func (p ExtParams) GetOutputBufferSizeOrThis(this int) int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, this)
		}
	}
	return this
}
func (p ExtParams) GetFieldMappings() []FieldMapping {
	if mappings, ok := p[FieldMappingsKey]; ok {
		if m, ok := mappings.([]FieldMapping); ok {
			return m
		}
	}
	return nil
}
func (p ExtParams) GetFMInputFields() []Field {
	var fields []Field
	for _, fm := range p.GetFieldMappings() {
		fields = append(fields, fm.SourceField)
	}
	return fields
}
func (p ExtParams) GetFMOutputFields() []Field {
	var fields []Field
	for _, fm := range p.GetFieldMappings() {
		fields = append(fields, fm.TargetField)
	}
	return fields
}
func (p ExtParams) GetInputDatasetFields() []Field {
	if dataset := p.GetInputDataset(); dataset != nil {
		return dataset.Fields
	}
	return nil
}
func (p ExtParams) GetOutputDatasetFields() []Field {
	if dataset := p.GetOutputDataset(); dataset != nil {
		return dataset.Fields
	}
	return nil
}

// todo check nil
func (e ExtParams) GetSourceFieldNameByTargetFieldName(targetFieldName string) string {
	mappings := e.GetFieldMappings()
	for _, fm := range mappings {
		if fm.TargetField.Name == targetFieldName {
			return fm.SourceField.Name
		}
	}
	return ""
}

// todo check nil
func (e ExtParams) GetSourceFieldByTarget(f Field) (*Field, error) {
	for _, fm := range e.GetFieldMappings() {
		if fm.TargetField.Name == f.Name {
			return &fm.SourceField, nil
		}
	}
	return nil, fmt.Errorf("field %s not found", f.Name)
}
