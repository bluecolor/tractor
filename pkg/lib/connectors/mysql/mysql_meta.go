package mysql

import (
	"regexp"

	"github.com/bluecolor/tractor/pkg/lib/cat/meta"
)

func (m *MySQLConnector) FetchDatasets(pattern string) ([]meta.Dataset, error) {
	datasets := []meta.Dataset{}
	rows, err := m.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		if pattern != "" {
			if match, _ := regexp.MatchString(pattern, tableName); !match {
				continue
			}
		}
		fields, err := m.fetchFields(tableName)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, meta.Dataset{
			Name:   tableName,
			Fields: fields,
		})
	}

	return datasets, nil
}
func (m *MySQLConnector) fetchFields(table string) ([]meta.Field, error) {
	result, err := m.db.Query("SHOW COLUMNS FROM " + table)
	if err != nil {
		return nil, err
	}
	fields := []meta.Field{}
	for result.Next() {
		var c column = column{}
		if err := result.Scan(&c.Field, &c.Type, &c.Null, &c.Key, &c.Default, &c.Extra); err != nil {
			return nil, err
		}
		fields = append(fields, c.getField())
	}
	return fields, nil
}
