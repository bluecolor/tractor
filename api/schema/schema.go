package schema

import "reflect"

// FT ...
type FT string

// FieldType  ..
var FieldType = struct {
	Numeric   FT
	String    FT
	Date      FT
	Timestamp FT
	Float     FT
}{
	Numeric:   "numeric",
	String:    "string",
	Date:      "date",
	Timestamp: "timestamp",
	Float:     "float",
}

// DecimalSize ...
type DecimalSize struct {
	Precision int64
	Scale     int64
	Ok        bool
}

// Nullable ...
type Nullable struct {
	Nullable bool
	Ok       bool
}

// Length ...
type Length struct {
	Length int64
	Ok     bool
}

// Field ...
type Field struct {
	Name        string
	Type        reflect.Type
	DecimalSize DecimalSize
	Nullable    Nullable
	Length      Length
}

// DataStore ...
type DataStore struct {
	Name   string
	Fields []Field
}
