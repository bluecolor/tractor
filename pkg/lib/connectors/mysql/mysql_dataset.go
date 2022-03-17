package mysql

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

func (c *MySQLConnector) getFields(options map[string]interface{}) ([]*types.Field, error) {

	database := c.config.Database
	if db, ok := options["database"]; ok {
		database = db.(string)
	}
	if _, ok := options["table"]; !ok {
		return nil, errors.New("table is required")
	}
	table := options["table"].(string)
	result, err := c.db.Query("SELECT column_name, data_type from information_schema.columns where table_schema = ? and table_name = ?", database, table)
	if err != nil {
		return nil, err
	}
	fields := []*types.Field{}
	for result.Next() {
		field := MySQLField{}
		if err := result.Scan(
			&field.ColumnName,
			&field.DataType,
		); err != nil {
			return nil, err
		}
		fields = append(fields, field.ToField())
	}
	return fields, nil
}

func (c *MySQLConnector) GetDataset(options map[string]interface{}) (*types.Dataset, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}
	defer c.Close()
	if fields, err := c.getFields(options); err != nil {
		return nil, err
	} else {
		return &types.Dataset{
			Name:   options["table"].(string),
			Fields: fields,
		}, nil
	}
}
