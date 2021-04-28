package cassandra

import (
	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	gocql "github.com/gocql/gocql"
	"github.com/mitchellh/mapstructure"
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
	return nil
}

func (c *Cassandra) Count() (int, error) {
	return c.count()
}

func (c *Cassandra) Init(catalog *config.Catalog) error {
	return c.connect()
}

func init() {
	inputs.Add("cassandra", func(config map[string]interface{}) tractor.Input {
		cass := Cassandra{}
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &cass,
			TagName:  "yaml",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(config)

		return &cass
	})

}
