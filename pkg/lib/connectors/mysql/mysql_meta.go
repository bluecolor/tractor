package mysql

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

func (c *MySQLConnector) FindFields(options map[string]interface{}) ([]types.Field, error) {
	database := c.config.Database
	if db, ok := options["database"]; ok {
		database = db.(string)
	}
	if _, ok := options["table"]; !ok {
		return nil, errors.New("table is required")
	}
	table := options["table"].(string)
	result, err := c.db.Query("SHOW COLUMNS FROM " + database + "." + table)
	if err != nil {
		return nil, err
	}
	fields := []types.Field{}
	for result.Next() {
		field := MySQLField{}
		if err := result.Scan(
			&field.Name,
			&field.Type,
			&field.Null,
			&field.Key,
			&field.Default,
			&field.Extra,
		); err != nil {
			return nil, err
		}
		fields = append(fields, field.ToField())
	}
	return fields, nil
}
