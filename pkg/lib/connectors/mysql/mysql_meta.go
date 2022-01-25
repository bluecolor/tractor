package mysql

import (
	"database/sql"
	"regexp"

	"github.com/bluecolor/tractor/pkg/lib/meta"
)

type columnConfig struct {
	Null    sql.NullString `json:"null"`
	Key     sql.NullString `json:"key"`
	Default sql.NullString `json:"default"`
	Extra   sql.NullString `json:"extra"`
}
type column struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	columnConfig
}

func (c *column) getConfig() map[string]interface{} {
	config := map[string]interface{}{}
	if c.Null.Valid {
		config["null"] = c.Null.String
	}
	if c.Key.Valid {
		config["key"] = c.Key.String
	}
	if c.Default.Valid {
		config["default"] = c.Default.String
	}
	if c.Extra.Valid {
		config["extra"] = c.Extra.String
	}
	return config
}
func (c *column) getField() meta.Field {
	return meta.Field{
		Name:   c.Field,
		Type:   c.Type,
		Config: c.getConfig(),
	}
}

func (m *MySQLConnector) FindDatasets(pattern string) ([]meta.Dataset, error) {
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
