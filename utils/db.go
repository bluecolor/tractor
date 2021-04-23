package utils

import (
	"database/sql"
	"fmt"
)

func GetFields(query string, db *sql.DB) (fields []Field, err error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	cts, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for _, ct := range cts {
		f := Field{
			Name: ct.Name(),
		}
		if precision, scale, ok := ct.DecimalSize(); ok {
			f.Precision = precision
			f.Scale = scale
		}
		if nullable, ok := ct.Nullable(); ok {
			f.Nullable = nullable
		}
		if length, ok := ct.Length(); ok {
			f.Length = length
		}
		fields = append(fields, f)
	}

	return fields, nil
}

func GetCount(query string, db *sql.DB) (count int, err error) {
	var q string = fmt.Sprintf("select count(1) from (%s)", query)
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		return count, err
	}
	rows.Next()
	if err := rows.Scan(count); err != nil {
		return count, err
	}
	return count, nil
}
