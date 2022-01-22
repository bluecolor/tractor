package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/bluecolor/tractor/pkg/models"
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

func (c *column) getConfig() []byte {
	config, err := json.Marshal(c.columnConfig)
	if err != nil {
		return nil
	}
	return config
}
func (c *column) getField() models.Field {
	return models.Field{
		Name:   c.Field,
		Type:   c.Type,
		Config: c.getConfig(),
	}
}

func NewMySQLConnector(config json.RawMessage) (*MySQLConnector, error) {
	var err error
	mysqlConfig := MySQLConfig{}
	if err = json.Unmarshal([]byte(config), &mysqlConfig); err != nil {
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
func (m *MySQLConnector) FetchDatasets() ([]models.Dataset, error) {
	return nil, nil
}
func (m *MySQLConnector) FetchDatasetsWithPattern(pattern string) ([]models.Dataset, error) {
	datasets := []models.Dataset{}
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
		datasets = append(datasets, models.Dataset{
			Name:   tableName,
			Fields: fields,
		})
	}

	return datasets, nil
}

func (m *MySQLConnector) fetchFields(table string) ([]models.Field, error) {
	result, err := m.db.Query("SHOW COLUMNS FROM " + table)
	if err != nil {
		return nil, err
	}
	fields := []models.Field{}
	for result.Next() {
		var c column = column{}
		if err := result.Scan(&c.Field, &c.Type, &c.Null, &c.Key, &c.Default, &c.Extra); err != nil {
			return nil, err
		}
		fields = append(fields, c.getField())
	}
	return fields, nil
}
