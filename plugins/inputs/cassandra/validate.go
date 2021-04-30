package cassandra

import "errors"

func (c *Cassandra) Validate() error {
	switch {
	case c.Cluster == "":
		return errors.New("Missing cluster")
	case c.Keyspace == "":
		return errors.New("Missing keyspace")
	case c.Table == "":
		return errors.New("Missing table")
	}

	return nil
}
