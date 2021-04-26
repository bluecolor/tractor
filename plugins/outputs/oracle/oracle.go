package oracle

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/outputs"
	"github.com/mitchellh/mapstructure"
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

	db *sql.DB
}

var catalogGiven bool = false
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

    catalog: see catalog messages //todo
`

func (o *Oracle) Description() string {
	return "Write to Oracle"
}

func (o *Oracle) SampleConfig() string {
	return sampleConfig
}

func (o *Oracle) Write(wire tractor.Wire) (err error) {

	tx, err := o.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			sendErrorFeed(wire, err)
		} else {
			tx.Commit()
			o.db.Close()
			wire.SendMessage(tractor.NewSuccessFeed(tractor.OutputPlugin))
		}
	}()

	for message := range wire.ReadMessages() {
		switch message.Type {
		case tractor.DataMessage:
			if data, ok := message.Content.(tractor.Data); ok {
				insertQuery, err = o.buildInsertQuery(len(data[0]))
				if err != nil {
					return err
				}
				err = insert(wire, tx, insertQuery, data)
				if err != nil {
					return err
				}
			} else {
				err = errors.New("Failed to cast data message")
				return err
			}
		case tractor.CatalogMessage:
			if !catalogGiven && strings.ToLower(o.Mode) == "drop-create" {
				if catalog, ok := message.Content.(config.Catalog); ok {
					err := o.dropCreate(&catalog)
					if err != nil {
						return err
					}
					if insertQuery == "" {
						insertQuery, err = o.buildInsertQuery(len(catalog.Properties))
						if err != nil {
							return err
						}
					}
				} else {
					err := errors.New("Failed to cast catalog message")
					return err
				}
			}
		}
	}
	return nil
}

func (o *Oracle) Init(catalog *config.Catalog) (err error) {
	if catalog != nil {
		catalogGiven = true
		if insertQuery == "" {
			insertQuery, err = o.buildInsertQuery(len(catalog.Properties))
			if err != nil {
				return err
			}
		}
		if strings.ToLower(o.Mode) == "drop-create" {
			err := o.dropCreate(catalog)
			if err != nil {
				return err
			}
		}
	}
	dsn, err := o.getDataSourceName()
	if err != nil {
		return err
	}
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return err
	}
	o.db = db
	return nil
}

func init() {
	outputs.Add("oracle", func(config map[string]interface{}) tractor.Output {
		oracle := Oracle{
			Port:      1521,
			Parallel:  1,
			BatchSize: 1000,
		}
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &oracle,
			TagName:  "yaml",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(config)

		return &oracle
	})
}
