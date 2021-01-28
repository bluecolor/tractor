package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/helpers/sqlhelper"
	"github.com/bluecolor/tractor/logging"
	"github.com/godror/godror"
	_ "github.com/godror/godror"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
	FetchSize        int    `yaml:"fetch_size"`
}

func (c *config) GetFetchSize() int {
	if c.FetchSize > 0 {
		return c.FetchSize
	}
	return 1000
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
func Run(wg *sync.WaitGroup, conf []byte, wire *api.Wire) error {

	config := config{}

	if err := api.ParseConfig(conf, &config); err != nil {
		return err
	}

	db, err := sql.Open("godror", getDataSourceName(config))

	if err != nil {
		return err
	}
	options := sqlhelper.Options{
		Db:           db,
		Query:        config.getQuery(),
		DbQueryArgs:  []interface{}{godror.FetchArraySize(config.GetFetchSize())},
		BatchSize:    config.GetFetchSize(),
		SendMetadata: true,
	}
	if err := sqlhelper.Send(wire, &options); err != nil {
		logging.Error(err)
		return err
	}

	wire.CloseMetadataChannel()
	wire.CloseDataChannel()
	db.Close()
	wg.Done()
	return nil
}

func getDataSourceName(config config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}

func (c *config) getQuery() string {
	return fmt.Sprintf("select * from %s", c.Table)
}
