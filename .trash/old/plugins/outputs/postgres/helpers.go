package postgres

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
		return fmt.Sprintf("%s text", name), nil
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
		return fmt.Sprintf("%s numeric(%d, %d)", name, precision, scale), nil
	}

	length := prop.Length
	if length == 0 {
		length = 4000
	}
	return fmt.Sprintf("%s text", sanitizeColumnName(prop.Name)), nil
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

func (o *Postgres) dropCreate(catalog *config.Catalog) error {
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

func (o *Postgres) buildInsertQuery(fieldCount int) (string, error) {
	if fieldCount == 0 {
		return "", errors.New("Field count is zero")
	}
	var columns = make([]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		columns[i] = "$" + strconv.Itoa(i+1)
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

func (p *Postgres) connect() error {
	info := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.Username, p.Password, p.Database)

	db, err := sql.Open("postgres", info)
	if err != nil {
		return err
	}
	p.db = db
	return nil
}
