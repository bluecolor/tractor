package mysql

import (
	"database/sql"
	"fmt"

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

func NewMySQLConnector(config connectors.ConnectorConfig) (connectors.Connector, error) {
	mysqlConfig := MySQLConfig{}
	if err := config.LoadConfig(&mysqlConfig); err != nil {
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
