package cassandra

import (
	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/utils"
	gocql "github.com/gocql/gocql"
)

type Cassandra struct {
	Cluster  string `yaml:"cluster"`
	Keyspace string `yaml:"keyspace"`
	Table    string `yaml:"table"`

	session *gocql.Session
}

var sampleConfig = `
    cluster: comma seperated list of hosts
    keyspace: name of the keyspace
    table: table name
`

func (c *Cassandra) Description() string {
	return "Read from Cassandra"
}

func (c *Cassandra) SampleConfig() string {
	return sampleConfig
}

func (c *Cassandra) Read(wire tractor.Wire) error {
	defer c.session.Close()
	return c.read(wire)
}

func (c *Cassandra) Count() (int, error) {
	return c.count()
}

func (c *Cassandra) Init(catalog *config.Catalog, params map[string]interface{}) error {
	return c.connect()
}

func newCassandra(options map[string]interface{}) *Cassandra {
	cass := &Cassandra{}
	utils.ParseOptions(cass, options)
	return cass
}

func init() {
	inputs.Add("cassandra", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Input, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newCassandra(options), nil
	})

}
