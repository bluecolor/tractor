package metadata

import (
	"database/sql"

	"github.com/bluecolor/tractor/api"
)

// GetFields ...
func GetFields(rows *sql.Rows) ([]api.Field, error) {
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var fields []api.Field
	for _, column := range columns {
		precision, scale, ok := column.DecimalSize()
		decimalSize := api.DecimalSize{Precision: precision, Scale: scale, Ok: ok}
		n, ok := column.Nullable()
		nullable := api.Nullable{Nullable: n, Ok: ok}
		l, ok := column.Length()
		length := api.Length{Length: l, Ok: ok}

		fields = append(fields, api.Field{
			Name:        column.Name(),
			Type:        column.ScanType(),
			DecimalSize: decimalSize,
			Nullable:    nullable,
			Length:      length,
		})
	}
	return fields, nil
}

// GetFieldsMetadata ...
func GetFieldsMetadata(rows *sql.Rows) (*api.Metadata, error) {

	fields, err := GetFields(rows)
	if err != nil {
		return nil, err
	}

	return &api.Metadata{
		Type:    api.FieldsMetadata,
		Content: fields,
	}, nil
}
