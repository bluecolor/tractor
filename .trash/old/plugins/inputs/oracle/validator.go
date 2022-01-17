package oracle

import "errors"

func (o *Oracle) Validate() error {
	switch {
	case o.Libdir == "":
		return errors.New("Missing libdir! Provide path to oracle instant client.")
	case o.URL == "" && o.Host == "":
		return errors.New("Missing host name! Provide IP or host name or url.")
	case o.URL == "" && o.Port == 0:
		return errors.New("Missing port number! Provide port no or url.")
	case o.URL == "" && o.Database == "":
		return errors.New("Missing database! Provide database name or url.")
	case o.Username == "":
		return errors.New("Missing username! Provide username.")
	case o.Password == "":
		return errors.New("Missing password! Provide password.")
	case o.Query == "" && o.Table == "":
		return errors.New("Missing source! Provide query or table.")
	}

	return nil
}
