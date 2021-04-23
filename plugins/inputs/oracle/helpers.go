package oracle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/utils"
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

func (o *Oracle) getFields() ([]utils.Field, error) {
	query, err := o.getQuery()
	if err != nil {
		return nil, err
	}
	return utils.GetFields(query, o.db)
}

func (o *Oracle) getQueries(parallel int) ([]string, error) {
	var colnames []string
	fields, err := o.getFields()
	if err != nil {
		return nil, err
	}
	query, err := o.getQuery()
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		colnames = append(colnames, f.Name)
	}
	count, err := utils.GetCount(query, o.db)
	if err != nil {
		return nil, err
	}
	chunkSize := (count / o.Parallel)
	q := fmt.Sprintf(
		"select * from (select * from (%s) order by %s)", query,
		fmt.Sprintf("order by %s", strings.Join(colnames, ",")),
	)
	queries := make([]string, o.Parallel)
	for i := 0; i < parallel; i++ {
		if i != parallel-1 {
			queries = append(queries, fmt.Sprintf(
				"%s where rownumber >= %d and rownumber < %d", q, i*chunkSize, (i+1)*chunkSize))
		} else {
			queries = append(queries, fmt.Sprintf(
				"%s where rownumber >= %d", q, i*chunkSize))
		}

	}
	return queries, nil
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
