package cassandra

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor"
	gocql "github.com/gocql/gocql"
)

func (c *Cassandra) getQuery() (string, error) {
	return fmt.Sprintf("select * from %s", c.Table), nil
}

func (c *Cassandra) connect() (err error) {
	cluster := gocql.NewCluster(strings.Split(c.Cluster, ",")...)
	cluster.Keyspace = c.Keyspace
	c.session, err = cluster.CreateSession()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cassandra) count() (count int, err error) {
	query := fmt.Sprintf("select count(1) as c from %s", c.Table)
	iter := c.session.Query(query).Iter()
	if ok := iter.Scan(count); !ok {
		return count, errors.New("Failed to get count")
	}
	return count, err
}

func (c *Cassandra) read(wire tractor.Wire) error {
	query, err := c.getQuery()
	if err != nil {
		return err
	}

	iter := c.session.Query(query).Iter()
	columns := iter.Columns()

	var record = make([]interface{}, len(columns))
	for i := range record {
		var v interface{}
		record[i] = &v
	}

	scanner := iter.Scanner()
	var data tractor.Data
	for scanner.Next() {
		if err := scanner.Scan(record...); err != nil {
			return err
		}
		data = append(data, record)
		if len(data) >= 100 { //todo
			wire.SendFeed(tractor.NewReadProgress(len(data)))
			wire.SendData(data)
			data = nil
		}
	}
	if len(data) > 0 {
		wire.SendData(data)
		data = nil
	}
	wire.SendFeed(tractor.NewSuccessFeed(tractor.InputPlugin))
	return err
}
