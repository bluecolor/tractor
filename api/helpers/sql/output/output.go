package ouput

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/logging"
)

// Options ...
type Options struct {
	Db            *sql.DB
	InsertQuery   string
	CreateQuery   string
	DbQueryArgs   []interface{}
	BatchSize     int
	Schema        string
	Table         string
	Mode          string
	Columns       []api.Field
	CreateOptions string
	TypeMap       map[reflect.Type]string
	BindPrefix    string
	Parallel      int
}

func (o *Options) getName(args ...string) string {
	var name string
	if len(args) > 0 {
		name = args[0]
	} else {
		name = o.Table
	}
	if o.Schema != "" {
		name = fmt.Sprintf("%s%s", o.Schema+".", name)
	}
	return name
}

// BuildCreateQuery ...
func (o *Options) BuildCreateQuery(args ...string) (string, error) {
	name := o.getName(args...)

	if len(o.Columns) == 0 {
		return "", errors.New("Columns not given")
	}

	columns := []string{}
	for _, c := range o.Columns {
		t, ok := o.TypeMap[c.Type]
		if !ok {
			return "", errors.New("Unknown field type for " + c.Name)
		}
		columns = append(columns, fmt.Sprintf("%s %s", c.Name, t))
	}
	var query = fmt.Sprintf(`create table %s (
        %s
    ) %s`, name, strings.Join(columns, ","), o.CreateOptions)
	return query, nil
}

// BuildInsertQuery ...
func (o *Options) BuildInsertQuery() (string, error) {
	if len(o.Columns) == 0 {
		return "", errors.New("Columns not given")
	}
	columns, values := []string{}, []string{}

	for i, c := range o.Columns {
		columns = append(columns, c.Name)
		values = append(values, o.BindPrefix+strconv.Itoa(i))
	}
	return fmt.Sprintf(
		"insert into %s (%s) values(%s)",
		o.Table, strings.Join(columns, ","), strings.Join(values, ","),
	), nil
}

// BuildDropQuery ...
func (o *Options) BuildDropQuery(args ...string) (string, error) {
	if o.Table == "" {
		return "", errors.New("Table name is not given")
	}
	return fmt.Sprintf("drop table %s", o.Table), nil
}

// BuildTruncateQuery ...
func (o *Options) BuildTruncateQuery() (string, error) {
	if o.Table == "" {
		return "", errors.New("Table name is not given")
	}
	return fmt.Sprintf("truncate table %s", o.Table), nil
}

// DropTable ...
func (h *Helper) DropTable() error {
	_, err := h.Db.Exec(h.BuildDropQuery())
	if err != nil {
		return err
	}
	return nil
}

// CreateTable ...
func (h *Helper) CreateTable() error {
	_, err := h.Db.Exec(h.BuildCreateQuery())
	if err != nil {
		return err
	}
	return nil
}

// TruncateTable ...
func (h *Helper) TruncateTable() error {
	_, err := h.Db.Exec(h.BuildTruncateQuery())
	if err != nil {
		return err
	}
	return nil
}

// Helper ...
type Helper struct {
	*Options
}

// NewHelper ...
func NewHelper(o *Options) *Helper {
	return &Helper{o}
}

// Run ...
func (h *Helper) Run(wire *api.Wire) error {

	if strings.ToLower(h.Mode) == "create" {
		if err := h.DropTable(); err != nil {
			logging.Warn("Failed to drop table", err)
		}
		if err := h.CreateTable(); err != nil {
			logging.Error("Failed to create target table")
			return err
		}
	} else if strings.ToLower(h.Mode) == "truncate" {
		if err := h.TruncateTable(); err != nil {
			logging.Error("Failed to truncate target table")
			return err
		}
	}

	insertQuery, err := h.BuildInsertQuery()
	if err != nil {
		logging.Error("Failed to build insert query")
		return err
	}
	tx, err := h.Db.Begin()
	if err != nil {
		logging.Error("Failed to start transaction")
		return err
	}

	for data := range wire.Data {
		for _, d := range data.Content {
			_, err = tx.Exec(insertQuery, d...)
			if err != nil {
				logging.Error("Failed to insert record")
				tx.Rollback()
				return err
			}
			wire.Feed <- api.NewWriteCountFeed(1)
		}
	}

	return nil
}

func (h *Helper) run(table string) error {

	return nil
}
