package types

type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeInt      FieldType = "int"
	FieldTypeNumber   FieldType = "number"
	FieldTypeBool     FieldType = "bool"
	FieldTypeDate     FieldType = "date"
	FieldTypeTime     FieldType = "time"
	FieldTypeDateTime FieldType = "datetime"
	FieldTypeArray    FieldType = "array"
	FieldTypeObject   FieldType = "object"
	FieldTypeMap      FieldType = "map"
	FieldTypeAny      FieldType = "any"
)

func (ftp FieldType) String() string {
	return string(ftp)
}
func FieldTypeFromString(s string) FieldType {
	switch s {
	case "string":
		return FieldTypeString
	case "int":
		return FieldTypeInt
	case "number":
		return FieldTypeNumber
	case "bool":
		return FieldTypeBool
	case "date":
		return FieldTypeDate
	case "time":
		return FieldTypeTime
	case "datetime":
		return FieldTypeDateTime
	case "array":
		return FieldTypeArray
	case "object":
		return FieldTypeObject
	case "map":
		return FieldTypeMap
	case "any":
		return FieldTypeAny
	}
	return FieldTypeString
}

type JsonSchema struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
}
type Dataset struct {
	Name   string   `json:"name"`
	Config Config   `json:"config"`
	Fields []*Field `json:"fields"`
}
type Field struct {
	Name   string    `json:"name"`
	Type   FieldType `json:"type"`
	Config Config    `json:"config"`
}

func (f *Field) GetExpressionOrName() string {
	if f.Config.GetString("expression", "") != "" {
		return f.Config.GetString("expression", "")
	}
	return f.Name
}
func (f *Field) GetOriginalType() string {
	return f.Config.GetString("originalType", "")
}

type FieldMapping struct {
	SourceField *Field `json:"sourceField"`
	TargetField *Field `json:"targetField"`
	Config      Config `json:"config"`
}

type Connection struct {
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Config         Config `json:"config"`
	AsSource       bool   `json:"asSource"`
	AsTarget       bool   `json:"asTarget"`
}
