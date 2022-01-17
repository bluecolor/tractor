package oracle

import (
	"database/sql"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/utils"
	"github.com/bluecolor/tractor/utils/db"
	_ "github.com/godror/godror"
)

type Oracle struct {
	Libdir    string `yaml:"libdir"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	URL       string `yaml:"url"`
	Query     string `yaml:"query"`
	Select    string `yaml:"select"`
	Where     string `yaml:"where"`
	Table     string `yaml:"table"`
	FetchSize int    `yaml:"fetch_size"`
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

  ## query to execute to fetch data
  query: select <fileds> from <table> where <conditions>

  select: comma seperated fields
  where: condition, filters
  ## if not schema name is given defaults to connection username
  table: table name with or without schema name. eg. "my_schema.my_table" or "my_table"

  fetch_size: defaults to 1000
  parallel: defaults to 1
`

func (o *Oracle) Description() string {
	return "Read from Oracle"
}

func (o *Oracle) SampleConfig() string {
	return sampleConfig
}

func (o *Oracle) Read(wire tractor.Wire) error {
	defer o.db.Close()
	queries, err := o.getQueries()
	if err != nil {
		println(err.Error())
		return err
	}
	for _, q := range queries {
		if err := db.Read(wire, q, o.db); err != nil {
			println(err.Error())
			return err
		}
	}
	return nil
}

func (o *Oracle) Init() error {
	return o.connect()
}

func (o *Oracle) Count() (int, error) {
	return o.count()
}

func newOracle(options map[string]interface{}) *Oracle {
	oracle := &Oracle{
		Port:      1521,
		Parallel:  1,
		FetchSize: 1000,
	}
	utils.ParseOptions(oracle, options)
	return oracle
}

func init() {
	inputs.Add("oracle", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Input, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newOracle(options), nil
	})
}