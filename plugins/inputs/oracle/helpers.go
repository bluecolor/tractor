package oracle

import (
	"errors"
	"fmt"
)

func (o *Oracle) getDataSourceName() (string, error) {

	if o.Username == "" || o.Password == "" {
		return "", errors.New("Missing credentials")
	}

	if o.URL != "" {
		return fmt.Sprintf(
			`user="%s" password="%s" connectString="%s"`,
			o.Username, o.Password, o.URL,
		), nil
	}
	if o.Host == "" || o.Port == 0 || o.Database == "" {
		return "", errors.New("Missing one or more connection information.(host, port, database)")
	}

	return fmt.Sprintf(
		`user="%s" password="%s" connectString="%s:%d/%s"`,
		o.Username, o.Password, o.Host, o.Port, o.Database,
	), nil
}

func (o *Oracle) getQuery() (string, error) {
	if o.Query != "" {
		return o.Query, nil
	}

	if o.Table == "" {
		return "", errors.New("Missing source table")
	}
	columns := "*"
	where := ""
	if o.Select != "" {
		columns = o.Select
	}
	if o.Where != "" {
		where = fmt.Sprintf("where %s", o.Where)
	}
	return fmt.Sprintf("select %s from %s %s", columns, o.Table, where), nil
}
