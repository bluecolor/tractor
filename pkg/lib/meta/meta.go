package meta

type Dataset struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	Fields []Field                `json:"fields"`
}

type Field struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type ExtInput struct {
	Dataset
	Parallel int `json:"parallel"`
}

type ExtOutput struct {
	Dataset
	Parallel      int            `json:"parallel"`
	FieldMappings []FieldMapping `json:"fieldMappings"`
}

type FieldMapping struct {
	SourceField Field                  `json:"sourceField"`
	TargetField Field                  `json:"targetField"`
	Config      map[string]interface{} `json:"config"`
}
