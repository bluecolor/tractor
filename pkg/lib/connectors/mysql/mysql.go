package mysql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm/utils"
)

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
type MySQLConnector struct {
	connectors.Base
	config MySQLConfig
	db     *sql.DB
}

func (m *MySQLConnector) GetDB() *sql.DB {
	return m.db
}

func New(config connectors.ConnectorConfig) (connectors.Connector, error) {
	mysqlConfig := MySQLConfig{}
	if err := config.LoadConfig(&mysqlConfig); err != nil {
		return nil, err
	}
	return &MySQLConnector{
		config: mysqlConfig,
	}, nil
}
func (c *MySQLConnector) Connect() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.config.Username, c.config.Password, c.config.Host, c.config.Port, c.config.Database,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// test connection
	if err := db.Ping(); err != nil {
		return err
	}
	c.db = db
	return nil
}
func (c *MySQLConnector) Close() error {
	return c.db.Close()
}
func (c *MySQLConnector) Validate(config connectors.ConnectorConfig) error {
	fields := reflect.VisibleFields(reflect.TypeOf(c.config))
	tags := make([]string, len(fields))
	for i, field := range fields {
		tags[i] = field.Tag.Get("json")
	}
	for key, _ := range config {
		if !utils.Contains(tags, key) {
			return fmt.Errorf("invalid config key %s", key)
		}
	}
	return nil
}

func init() {
	connectors.Add("mysql", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config)
	})
}
