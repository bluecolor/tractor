package postgres

import (
	"database/sql"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/outputs"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
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

	db *sql.DB
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

func (p *Postgres) Write(wire tractor.Wire) (err error) {
	defer p.db.Close()
	return nil
}

func (p *Postgres) Init(catalog *config.Catalog) (err error) {
	err = p.connect()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	outputs.Add("postgres", func(config map[string]interface{}) tractor.Output {
		pg := Postgres{
			Port:      1521,
			Parallel:  1,
			BatchSize: 1000,
		}
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &pg,
			TagName:  "yaml",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(config)

		return &pg
	})
}
