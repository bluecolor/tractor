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
