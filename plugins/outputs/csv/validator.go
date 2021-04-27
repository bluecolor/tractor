package csv

import "errors"

func (o *Csv) ValidateConfig() error {
	switch {
	case o.Path == "":
		return errors.New("Missing path")
	}

	return nil
}
