package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bluecolor/tractor/api"
	_ "github.com/godror/godror"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
}

// SampleConfig ...
func SampleConfig() string {
	return `
plugin: oracle
username: username
password: password
libdir: path/to/oracle_instantclient
connection_string: "localhost:1521/orcl"
table: my_table
schema:
    name: datastore_name
    fileds:
        - name: filed_name
          type: data_type
          precision: precision
          scale: scale
`
}

// Description ...
func Description() string {
	return "Read data from oracle database"
}

// PluginType ...
func PluginType() api.PluginType {
	return api.InputPlugin
}

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, channel chan []byte) {
	config, err := getConfig(conf)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("godror", getDataSourceName(config))
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(config.getQuery())
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		data := ""
		if err := rows.Scan(&data); err != nil {
			panic(err)
		}
		println(data)
	}
	db.Close()
	wg.Done()
}

func getConfig(conf []byte) (*config, error) {
	config := config{}
	if err := api.ParseConfig(conf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getDataSourceName(config *config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}

func (c *config) getQuery() string {
	return fmt.Sprintf("select * from %s", c.Table)
}
