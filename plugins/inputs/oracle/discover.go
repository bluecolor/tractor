package oracle

import (
	"database/sql"

	"github.com/bluecolor/tractor/config"
)

func (o *Oracle) Discover() (*config.Catalog, error) {

	query, err := o.getQuery()

	if err != nil {
		return nil, err
	}
	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}

	ctypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	catalog := config.Catalog{
		Name: o.Table,
	}

	for _, ct := range ctypes {
		prop := getPorperty(ct)
		catalog.Properties = append(catalog.Properties, prop)
	}

	return &catalog, nil
}

func getPorperty(ct *sql.ColumnType) config.Property {
	prop := config.Property{}

	switch ct.DatabaseTypeName() {
	case "VARCHAR2", "NVARCHAR2", "CHAR", "LONG", "VARCHAR", "NCHAR":
		prop.Type = "string"
		if length, ok := ct.Length(); ok {
			prop.Length = length
		}
		return prop
	case "NUMBER":
		prop.Type = "numeric"
		if precision, scale, ok := ct.DecimalSize(); ok {
			prop.Precision = precision
			prop.Scale = scale
		}
		if length, ok := ct.Length(); ok {
			prop.Length = length
		}
		return prop
	case "DATE", "TIMESTAMP":
		prop.Type = "string"
		if precision, scale, ok := ct.DecimalSize(); ok {
			prop.Precision = precision
			prop.Scale = scale
		}
		if length, ok := ct.Length(); ok {
			prop.Length = length
		}
		return prop
	}

	if precision, scale, ok := ct.DecimalSize(); ok {
		prop.Precision = precision
		prop.Scale = scale
	}
	if length, ok := ct.Length(); ok {
		prop.Length = length
	}
	return prop
}
