package md

import (
	"reflect"

	"github.com/bluecolor/tractor/api/md/mdt"
)

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

// Metadata ...
type Metadata struct {
	Type    mdt.Type
	Content interface{}
}
