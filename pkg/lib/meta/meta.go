package meta

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
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config Config `json:"config"`
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
	SourceField Field  `json:"sourceField"`
	TargetField Field  `json:"targetField"`
	Config      Config `json:"config"`
}
