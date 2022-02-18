package mysql

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
)

func (c *MySQLConnector) FindFields(options map[string]interface{}) ([]params.Field, error) {
	database := c.config.Database
	if db, ok := options["database"]; ok {
		database = db.(string)
	}
	table := options["table"].(string)
	result, err := c.db.Query("SHOW COLUMNS FROM " + database + "." + table)
	if err != nil {
		return nil, err
	}
	fields := []params.Field{}
	for result.Next() {
		var name, tp string
		if err := result.Scan(&name, &tp); err != nil {
			return nil, err
		}
		fields = append(fields, params.Field{
			Name: name,
			Type: tp,
		})
	}
	return fields, nil
}
