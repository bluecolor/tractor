package oracle

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	dbu "github.com/bluecolor/tractor/utils/db"
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

func sanitizeColumnName(column string, args ...string) string {
	replace := "_"
	if len(args) > 0 {
		replace = args[0]
	}
	r := strings.NewReplacer(" ", replace, "%", replace)
	return r.Replace(column)
}

func columnFromProperty(prop *config.Property) (string, error) {
	if prop.Name == "" {
		return "", errors.New("Missing property name")
	}
	name := sanitizeColumnName(prop.Name)
	switch prop.Type {
	case "string":
		length := prop.Length
		if length == 0 {
			length = 4000
		}
		return fmt.Sprintf("%s varchar2(%d)", name, length), nil
	case "date":
		return fmt.Sprintf("%s timestamp", name), nil
	case "numeric":
		precision := prop.Precision
		scale := prop.Scale
		if precision == 0 {
			precision = 22
		}
		if scale >= 22 {
			scale = 21
		}
		return fmt.Sprintf("%s number(%d, %d)", name, precision, scale), nil
	}

	length := prop.Length
	if length == 0 {
		length = 4000
	}
	return fmt.Sprintf("%s varchar2(%d)", sanitizeColumnName(prop.Name), length), nil
}

func columnsFromProperties(properties []config.Property) (columns []string, err error) {
	for _, p := range properties {
		column, err := columnFromProperty(&p)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

func (o *Oracle) dropCreate(catalog *config.Catalog) error {
	table := o.Table
	if table == "" {
		table = catalog.Name
	}
	if table == "" {
		return errors.New("Table name is missing")
	}
	_ = dbu.DropTable(o.db, table)

	columns, err := columnsFromProperties(catalog.Properties)
	if err != nil {
		return err
	}

	return dbu.CreateTable(o.db, table, columns, "")
}

func (o *Oracle) buildInsertQuery(fieldCount int) (string, error) {
	if fieldCount == 0 {
		return "", errors.New("Field count is zero")
	}
	var columns = make([]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		columns[i] = ":" + strconv.Itoa(i)
	}
	return fmt.Sprintf("insert into %s values(%s)", o.Table, strings.Join(columns, ", ")), nil
}

func sendErrorFeed(wire tractor.Wire, err error) error {
	feed := tractor.NewErrorFeed(tractor.OutputPlugin, err)
	wire.SendFeed(feed)
	return err
}

func insert(wire tractor.Wire, tx *sql.Tx, query string, data tractor.Data) error {
	count, err := dbu.Insert(tx, query, data)
	if err != nil {
		return err
	}
	progress := tractor.NewWriteProgress(count)
	wire.SendFeed(progress)
	return nil
}

func (o *Oracle) connect() error {
	dsn, err := o.getDataSourceName()
	if err != nil {
		return err
	}
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return err
	}
	o.db = db
	return nil
}
