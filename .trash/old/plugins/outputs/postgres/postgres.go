package postgres

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/outputs"
	"github.com/bluecolor/tractor/utils"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Mode      string `yaml:"mode"`
	Table     string `yaml:"table"`
	BatchSize int    `yaml:"batch_size"`
	Parallel  int    `yaml:"parallel"`

	catalog *config.Catalog
	db      *sql.DB
}

var insertQuery string = ""

var sampleConfig = `
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

func (p *Postgres) Description() string {
	return "Write to Postgres"
}

func (p *Postgres) SampleConfig() string {
	return sampleConfig
}

func (o *Postgres) startWorker(wire tractor.Wire) (err error) {
	tx, err := o.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			sendErrorFeed(wire, err)
		} else {
			tx.Commit()
			o.db.Close()
			wire.SendFeed(tractor.NewSuccessFeed(tractor.OutputPlugin))
		}
	}()

	for data := range wire.ReadData() {
		err = insert(wire, tx, insertQuery, data)
		if err != nil {
			println(err.Error())
			return err
		}
	}
	return nil
}

func (p *Postgres) Write(wire tractor.Wire) (err error) {
	defer p.db.Close()
	if insertQuery == "" {
		for data := range wire.ReadData() {
			insertQuery, err = p.buildInsertQuery(len((data)[0]))
			if err != nil {
				return err
			}
			wire.SendData(data)
			break
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < p.Parallel; i++ {
		go func(wg *sync.WaitGroup) {
			p.startWorker(wire)
			wg.Done()
		}(&wg)
		wg.Add(1)
	}
	wg.Wait()
	return nil

}

func (p *Postgres) Init(catalog *config.Catalog) (err error) {
	err = p.connect()
	if err != nil {
		return err
	}

	if catalog != nil {
		if insertQuery == "" {
			insertQuery, err = p.buildInsertQuery(len(catalog.Properties))
			if err != nil {
				return err
			}
		}
		if strings.ToLower(p.Mode) == "drop-create" {
			err := p.dropCreate(catalog)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newPostgres(options map[string]interface{}) *Postgres {
	oracle := &Postgres{
		Port:      1521,
		Parallel:  1,
		BatchSize: 1000,
	}
	utils.ParseOptions(oracle, options)
	return oracle
}

func init() {
	outputs.Add("postgres", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Output, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		oracle := newPostgres(options)
		oracle.catalog = catalog

		return oracle, nil
	})
}
