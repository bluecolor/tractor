package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	SessionIDKey      = "session_id"
	TimeoutKey        = "timeout"
	InputDatasetKey   = "input_dataset"
	OutputDatasetKey  = "output_dataset"
	ParallelKey       = "parallel"
	ExtractionModeKey = "extraction_mode"
	FieldMappingsKey  = "field_mappings"
	BufferSizeKey     = "buffer_size"
	DefaultTimeOut    = time.Second * 60 * 10 // todo from env 10 minutes
)

type SessionParams map[string]interface{}

func (p SessionParams) Init() SessionParams {
	return p.EnsureSessionID().WithTimeout(DefaultTimeOut)
}
func (p SessionParams) WithTimeout(timeout time.Duration) SessionParams {
	p[TimeoutKey] = timeout
	return p
}
func (p SessionParams) WithSessionID(id interface{}) SessionParams {
	switch val := id.(type) {
	case string:
		p[SessionIDKey] = val
	default:
		p[SessionIDKey] = fmt.Sprintf("%d", val)
	}
	return p
}
func (p SessionParams) GetSessionID() string {
	if id, ok := p[SessionIDKey]; ok {
		if i, ok := id.(string); ok {
			return i
		}
	}
	return ""
}
func (p SessionParams) EnsureSessionID() SessionParams {
	if _, ok := p[SessionIDKey]; !ok {
		p.WithSessionID(uuid.New().String())
	}
	return p
}
func (p SessionParams) WithInputDataset(dataset *Dataset) SessionParams {
	p[InputDatasetKey] = dataset
	return p
}
func (p SessionParams) WithInputParallel(parallel int) SessionParams {
	p[InputDatasetKey].(*Dataset).Config[ParallelKey] = parallel
	return p
}
func (p SessionParams) WithOutputParallel(parallel int) SessionParams {
	p[OutputDatasetKey].(*Dataset).Config[ParallelKey] = parallel
	return p
}
func (p SessionParams) WithOutputDataset(dataset *Dataset) SessionParams {
	p[OutputDatasetKey] = dataset
	return p
}
func (p SessionParams) WithExtractionModeString(mode string) SessionParams {
	p[ExtractionModeKey] = ExtractionModeFromString(mode)
	return p
}
func (p SessionParams) WithFieldMappings(mappings []FieldMapping) SessionParams {
	p[FieldMappingsKey] = mappings
	return p
}
func (p SessionParams) GetTimeout() time.Duration {
	if timeout, ok := p[TimeoutKey]; ok {
		if t, ok := timeout.(time.Duration); ok {
			return t
		}
	}
	return DefaultTimeOut
}
func (p SessionParams) GetExtractionMode() ExtractionMode {
	if mode, ok := p[ExtractionModeKey]; ok {
		if m, ok := mode.(ExtractionMode); ok {
			return m
		}
	}
	return ExtractionModeCreate
}
func (p SessionParams) GetInputDataset() *Dataset {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(*Dataset); ok {
			return d
		}
	}
	return nil
}
func (p SessionParams) GetOutputDataset() *Dataset {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(*Dataset); ok {
			return d
		}
	}
	return nil
}
func (p SessionParams) GetInputParallel() int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(ParallelKey, 1)
		}
	}
	return 1
}
func (p SessionParams) GetOutputParallel() int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(ParallelKey, 1)
		}
	}
	return 1
}
func (p SessionParams) GetInputBufferSize() int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, 1000) // todo default from env
		}
	}
	return 1000
}
func (p SessionParams) GetOutputBufferSize() int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, 1000) // todo default from env
		}
	}
	return 1000
}
func (p SessionParams) GetInputBufferSizeOrThis(this int) int {
	if dataset, ok := p[InputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, this)
		}
	}
	return this
}
func (p SessionParams) GetOutputBufferSizeOrThis(this int) int {
	if dataset, ok := p[OutputDatasetKey]; ok {
		if d, ok := dataset.(Dataset); ok {
			return d.Config.GetInt(BufferSizeKey, this)
		}
	}
	return this
}
func (p SessionParams) GetFieldMappings() []FieldMapping {
	if mappings, ok := p[FieldMappingsKey]; ok {
		if m, ok := mappings.([]FieldMapping); ok {
			return m
		}
	}
	return nil
}
func (p SessionParams) GetFMInputFields() []*Field {
	var fields []*Field
	for _, fm := range p.GetFieldMappings() {
		fields = append(fields, fm.SourceField)
	}
	return fields
}
func (p SessionParams) GetFMOutputFields() []*Field {
	var fields []*Field
	for _, fm := range p.GetFieldMappings() {
		fields = append(fields, fm.TargetField)
	}
	return fields
}
func (p SessionParams) GetInputDatasetFields() []*Field {
	if dataset := p.GetInputDataset(); dataset != nil {
		return dataset.Fields
	}
	return nil
}
func (p SessionParams) GetOutputDatasetFields() []*Field {
	if dataset := p.GetOutputDataset(); dataset != nil {
		return dataset.Fields
	}
	return nil
}

// todo check nil
func (p SessionParams) GetSourceFieldNameByTargetFieldName(targetFieldName string) string {
	mappings := p.GetFieldMappings()
	for _, fm := range mappings {
		if fm.TargetField.Name == targetFieldName {
			return fm.SourceField.Name
		}
	}
	return ""
}

// todo check nil
func (p SessionParams) GetSourceFieldByTarget(f Field) (*Field, error) {
	for _, fm := range p.GetFieldMappings() {
		if fm.TargetField.Name == f.Name {
			return fm.SourceField, nil
		}
	}
	return nil, fmt.Errorf("field %s not found", f.Name)
}
