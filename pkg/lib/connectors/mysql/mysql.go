package mysql

import (
	"database/sql"
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/cat/meta"
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
type MySQLConnector struct {
	config MySQLConfig
	db     *sql.DB
}

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

func NewMySQLConnector(config connectors.ConnectorConfig) (connectors.Connector, error) {
	var err error
	mysqlConfig := MySQLConfig{}
	if err = config.LoadConfig(&mysqlConfig); err != nil {
		return nil, err
	}
	return &MySQLConnector{
		config: mysqlConfig,
	}, nil
}
func (m *MySQLConnector) Connect() error {
	dsn := m.config.User + ":" + m.config.Password + "@tcp(" + m.config.Host + ":" + fmt.Sprint(m.config.Port) + ")/" + m.config.Database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}
func (m *MySQLConnector) Close() error {
	return m.db.Close()
}
func init() {
	connectors.Add("mysql", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return NewMySQLConnector(config)
	})
}
