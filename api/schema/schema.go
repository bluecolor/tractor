package schema

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

// Field ...
type Field struct {
	Name      string
	Type      FT
	Precision uint
	Scale     uint
}

// DataStore ...
type DataStore struct {
	Name   string
	Fields []Field
}
