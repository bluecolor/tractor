package meta

type Field struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type Dataset struct {
	Name   string            `json:"name"`
	Fields []Field           `json:"fields"`
	Config map[string]string `json:"config"`
}
