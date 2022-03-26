package mysql

import "errors"

func (c *MySQLConnector) GetInfo(info string, options map[string]interface{}) (interface{}, error) {
	switch info {
	case "databases":
		return c.getDatabases(options)
	case "tables":
		return c.getTables(options)
	}
	return nil, errors.New("unknown info type")
}

func (c *MySQLConnector) getDatabases(options map[string]interface{}) (interface{}, error) {
	query := "SHOW DATABASES"
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	return databases, nil
}
func (c *MySQLConnector) getTables(options map[string]interface{}) (interface{}, error) {
	database := options["database"].(string)
	query := "SHOW TABLES FROM " + database
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
