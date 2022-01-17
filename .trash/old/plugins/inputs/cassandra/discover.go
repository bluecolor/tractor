package cassandra

import (
	"github.com/bluecolor/tractor/config"
	gocql "github.com/gocql/gocql"
)

func (c *Cassandra) Discover() (*config.Catalog, error) {
	query, err := c.getQuery()
	if err != nil {
		return nil, err
	}
	iter := c.session.Query(query).Iter()
	columns := iter.Columns()
	catalog := &config.Catalog{
		Name: c.Table,
	}
	for _, column := range columns {
		property := getPorperty(column)
		catalog.Properties = append(catalog.Properties, property)
	}
	return catalog, nil
}

func getPorperty(c gocql.ColumnInfo) config.Property {
	prop := config.Property{
		Name: c.Name,
	}
	switch c.TypeInfo.Type() {
	case
		gocql.TypeBigInt, gocql.TypeTinyInt, gocql.TypeInt, gocql.TypeDecimal,
		gocql.TypeCounter, gocql.TypeFloat, gocql.TypeVarint:
		prop.Type = "numeric"
	case
		gocql.TypeDate, gocql.TypeTime, gocql.TypeTimestamp:
		prop.Type = "date"
	default:
		prop.Type = "string"
	}
	return prop
}
