package mysql

import "errors"

func (c *MySQLConnector) Resolve(request string, body map[string]interface{}) (interface{}, error) {
	resolvers := c.GetResolvers()
	if resolver, ok := resolvers[request]; ok {
		return func(body map[string]interface{}) (interface{}, error) {
			if err := c.Connect(); err != nil {
				return nil, err
			}
			defer c.Close()
			return resolver(body)
		}(body)
	} else {
		return nil, errors.New("unknown request")
	}
}
func (c *MySQLConnector) GetResolvers() map[string]func(map[string]interface{}) (interface{}, error) {
	return map[string]func(map[string]interface{}) (interface{}, error){
		"databases": c.ResolveDatabases,
		"tables":    c.ResolveTables,
	}
}
func (c *MySQLConnector) ResolveDatabases(options map[string]interface{}) (interface{}, error) {
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
func (c *MySQLConnector) ResolveTables(options map[string]interface{}) (interface{}, error) {
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
