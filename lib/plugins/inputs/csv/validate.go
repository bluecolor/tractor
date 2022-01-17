package csv

import "errors"

func (c *Csv) Validate() error {
	switch {
	case c.Path == "":
		return errors.New("missing path")
	}

	return nil
}
