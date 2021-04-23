package oracle

import (
	"database/sql"

	"github.com/bluecolor/tractor"
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
	URL       string `yaml:"url"`
	Table     string `yaml:"table"`
	BatchSize int    `yaml:"batch_size"`
	Parallel  int    `yaml:"parallel"`

	db *sql.DB
}

var sampleConfig = `
  ## instant client from oracle
  ## https://www.oracle.com/database/technologies/instant-client/downloads.html
  libdir: "/path/to/oracle_instant_client"

  host: "host name or ip address"
  port: port number default 1521
  database: service name or sid
  username: connection username
  password: connection password

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

func (o *Oracle) Write(ch <-chan *tractor.Message) error {

	for message := range ch {

	}

	return nil
}

func (o *Oracle) Init() error {
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
