package mysql

import (
	"database/sql"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type MySQLField struct {
	ColumnName             string
	OrdinalPosition        int
	ColumnDefault          sql.NullString
	IsNullable             string
	DataType               string
	CharacterMaximumLength int
	CharacterOctetLength   int
	NumericPrecision       int
	NumericScale           int
	DatetimePrecision      int
	CharacterSetName       sql.NullString
	ColumnType             string
	ColumnKey              sql.NullString
	Extra                  sql.NullString
	ColumnComment          sql.NullString
}

func (f MySQLField) ToField() *types.Field {
	typeMap := map[string]types.FieldType{
		"tinyint":    types.FieldTypeInt,
		"smallint":   types.FieldTypeInt,
		"mediumint":  types.FieldTypeInt,
		"int":        types.FieldTypeInt,
		"bigint":     types.FieldTypeInt,
		"float":      types.FieldTypeNumber,
		"double":     types.FieldTypeNumber,
		"decimal":    types.FieldTypeNumber,
		"bit":        types.FieldTypeInt,
		"bool":       types.FieldTypeBool,
		"char":       types.FieldTypeString,
		"varchar":    types.FieldTypeString,
		"tinytext":   types.FieldTypeString,
		"text":       types.FieldTypeString,
		"mediumtext": types.FieldTypeString,
		"longtext":   types.FieldTypeString,
		"date":       types.FieldTypeDate,
		"time":       types.FieldTypeTime,
		"datetime":   types.FieldTypeDateTime,
		"timestamp":  types.FieldTypeDateTime,
		"year":       types.FieldTypeInt,
		"enum":       types.FieldTypeString,
		"set":        types.FieldTypeString,
		"binary":     types.FieldTypeString,
		"varbinary":  types.FieldTypeString,
		"tinyblob":   types.FieldTypeString,
		"blob":       types.FieldTypeString,
		"mediumblob": types.FieldTypeString,
		"longblob":   types.FieldTypeString,
		"json":       types.FieldTypeObject,
	}

	tp, ok := typeMap[f.DataType]
	if !ok {
		tp = types.FieldTypeString
	}

	return &types.Field{
		Name:    f.ColumnName,
		Type:    tp,
		RawType: f.ColumnType,
		Config:  types.Config{},
	}
}
func ToMySQLField(f types.Field) MySQLField {
	typeMap := map[types.FieldType]string{
		types.FieldTypeInt:      "int",
		types.FieldTypeNumber:   "double",
		types.FieldTypeBool:     "bool",
		types.FieldTypeString:   "string",
		types.FieldTypeDate:     "date",
		types.FieldTypeTime:     "time",
		types.FieldTypeDateTime: "datetime",
		types.FieldTypeObject:   "json",
	}

	tp, ok := typeMap[f.Type]
	if !ok {
		tp = "string"
	}

	null := sql.NullString{
		String: "YES",
		Valid:  true,
	}
	if s, ok := f.Config["null"]; ok {
		null.Valid = true
		null.String = s.(string)
	}
	key := sql.NullString{
		String: "",
		Valid:  false,
	}
	if s, ok := f.Config["key"]; ok {
		key.Valid = true
		key.String = s.(string)
	}
	defaultValue := sql.NullString{
		String: "",
		Valid:  false,
	}
	if s, ok := f.Config["default"]; ok {
		defaultValue.Valid = true
		defaultValue.String = s.(string)
	}
	extra := sql.NullString{
		String: "",
		Valid:  false,
	}
	if s, ok := f.Config["extra"]; ok {
		extra.Valid = true
		extra.String = s.(string)
	}

	return MySQLField{
		ColumnName: f.Name,
		DataType:   tp,
	}
}
