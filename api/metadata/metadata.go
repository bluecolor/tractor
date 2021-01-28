package metadata

import (
	"reflect"
)

// Type ...
type Type int

const (
	// Fields ...
	Fields Type = iota
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

// Metadata ...
type Metadata struct {
	Type    Type
	Content interface{}
}

// Data ...
type Data struct {
	Content [][]interface{}
}
