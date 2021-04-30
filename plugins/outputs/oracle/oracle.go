package oracle

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/outputs"
	"github.com/bluecolor/tractor/utils"
)

type Oracle struct {
	Libdir    string `yaml:"libdir"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Mode      string `yaml:"mode"`
	URL       string `yaml:"url"`
	Table     string `yaml:"table"`
	BatchSize int    `yaml:"batch_size"`
	Parallel  int    `yaml:"parallel"`

	catalog *config.Catalog
	db      *sql.DB
}

var insertQuery string = ""

var sampleConfig = `
    ## instant client from oracle
    ## https://www.oracle.com/database/technologies/instant-client/downloads.html
    libdir: "/path/to/oracle_instant_client"

    host: "host name or ip address"
    port: port number default 1521
    database: service name or sid
    username: connection username
    password: connection password

    mode: TRUNCATE, DROP-CREATE, DEFAULT
    ## if not schema name is given defaults to connection username
    table: table name with or without schema name. eg. "my_schema.my_table" or "my_table"

    batch_size: defaults to 1000
    parallel: defaults to 1
`

func (o *Oracle) Description() string {
	return "Write to Oracle"
}

func (o *Oracle) SampleConfig() string {
	return sampleConfig
}

func (o *Oracle) startWorker(wire tractor.Wire) (err error) {
	tx, err := o.db.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback()
			sendErrorFeed(wire, err)
		} else {
			err = tx.Commit()
			err = o.db.Close()
			wire.SendFeed(tractor.NewSuccessFeed(tractor.OutputPlugin))
		}
	}()

	for data := range wire.ReadData() {
		err = insert(wire, tx, insertQuery, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Oracle) Write(wire tractor.Wire) (err error) {

	if insertQuery == "" {
		for data := range wire.ReadData() {
			insertQuery, err = o.buildInsertQuery(len((data)[0]))
			if err != nil {
				println(err.Error())
				return err
			}
			wire.SendData(data)
			break
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < o.Parallel; i++ {
		go func(wg *sync.WaitGroup) {
			o.startWorker(wire)
			wg.Done()
		}(&wg)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}

func (o *Oracle) Init() (err error) {
	err = o.connect()
	if err != nil {
		return err
	}

	if o.catalog != nil {
		if insertQuery == "" {
			insertQuery, err = o.buildInsertQuery(len(o.catalog.Properties))
			if err != nil {
				return err
			}
		}
		if strings.ToLower(o.Mode) == "drop-create" {
			err := o.dropCreate(o.catalog)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newOracle(options map[string]interface{}) *Oracle {
	oracle := &Oracle{
		Port:      1521,
		Parallel:  1,
		BatchSize: 1000,
	}
	utils.ParseOptions(oracle, options)
	return oracle
}

func init() {
	outputs.Add("oracle", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Output, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		oracle := newOracle(options)
		oracle.catalog = catalog

		return oracle, nil
	})
}
