package csv

import "errors"

func (c *Csv) Validate() error {
	switch {
	case c.Path == "":
		return errors.New("Missing path")
	}

	return nil
}
