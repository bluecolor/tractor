package csv

import "errors"

func (c *Csv) ValidateConfig() error {
	switch {
	case c.Path == "":
		return errors.New("Missing path")
	}

	return nil
}
