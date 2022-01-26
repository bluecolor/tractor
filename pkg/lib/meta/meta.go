package meta

import "fmt"

var (
	ExtractionModeCreate = "create"
	ExtractionModeInsert = "insert"
	ExtractionModeMerge  = "merge"
	ExtractionModeAppend = "append"
)

type Config map[string]interface{}

func (c Config) GetString(key string, def string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}
func (c Config) GetInt(key string, def int) int {
	if v, ok := c[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return def
}

type Dataset struct {
	Name   string  `json:"name"`
	Config Config  `json:"config"`
	Fields []Field `json:"fields"`
}
type Field struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Expression string `json:"expression"`
	Config     Config `json:"config"`
}

func (f *Field) GetExpressionOrName() string {
	if f.Expression != "" {
		return f.Expression
	}
	return f.Name
}

type ExtInput struct {
	Dataset
	Parallel int `json:"parallel"`
}
type FieldMapping struct {
	SourceField Field  `json:"sourceField"`
	TargetField Field  `json:"targetField"`
	Config      Config `json:"config"`
}
type ExtOutput struct {
	Dataset
	Parallel       int            `json:"parallel"`
	FieldMappings  []FieldMapping `json:"fieldMappings"`
	ExtractionMode string         `json:"extractionMode"`
}

func (e *ExtOutput) GetSourceFieldNameByTargetFieldName(targetFieldName string) string {
	for _, fm := range e.FieldMappings {
		if fm.TargetField.Name == targetFieldName {
			return fm.SourceField.Name
		}
	}
	return ""
}
func (e *ExtOutput) GetSourceFieldByTarget(f Field) (*Field, error) {
	for _, fm := range e.FieldMappings {
		if fm.TargetField.Name == f.Name {
			return &fm.SourceField, nil
		}
	}
	return nil, fmt.Errorf("field %s not found", f.Name)
}
