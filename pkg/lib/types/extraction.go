package types

import (
	"sort"
	"time"
)

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

type Connection struct {
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Config         Config `json:"config"`
}
type Dataset struct {
	Name       string      `json:"name"`
	Connection *Connection `json:"connection"`
	Config     Config      `json:"config"`
	Fields     []*Field    `json:"fields"`
}

func (d *Dataset) GetParallel() int {
	if d.Config == nil {
		return 1
	}
	return d.Config.GetInt("parallel", 1)
}
func (d *Dataset) GetBufferSize() int {
	if d.Config == nil {
		return 1000
	}
	return d.Config.GetInt("bufferSize", 1000)
}
func (d *Dataset) GetExtractionMode(def ...string) string {
	if d.Config == nil {
		return ""
	}
	if len(def) > 0 {
		return d.Config.GetString("extractionMode", def[0])
	}
	return d.Config.GetString("extractionMode", "")
}
func (d *Dataset) GetFields() []*Field {
	sort.Slice(d.Fields, func(i, j int) bool {
		return d.Fields[i].Order < d.Fields[j].Order
	})
	return d.Fields
}

type Field struct {
	Name    string    `json:"name"`
	Order   int       `json:"order"`
	Type    FieldType `json:"type"`
	RawType string    `json:"rawType"`
	Config  Config    `json:"config"`
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

type Extraction struct {
	Name          string   `json:"name"`
	SourceDataset *Dataset `json:"sourceDataset"`
	TargetDataset *Dataset `json:"targetDataset"`
	Config        Config   `json:"config"`
}

func (e *Extraction) GetTimeout() time.Duration {
	return time.Duration(e.Config.GetInt("timeout", 60*30)) * time.Second
}

type Session struct {
	ID         string      `json:"id"`
	Extraction *Extraction `json:"extraction"`
	Config     Config      `json:"config"`
}
