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
	"github.com/bluecolor/tractor/util"
)

type messageType int

const (
	errorMessage messageType = iota
	successMessage
	progressMessage
	stopOrder
)

type message struct {
	Type    messageType
	Sender  string
	Content interface{}
}

// Options ...
type Options struct {
	Db            *sql.DB
	InsertQuery   string
	CreateQuery   string
	DbQueryArgs   []interface{}
	BatchSize     int
	Schema        string
	TempSchema    string
	Table         string
	Mode          string
	Columns       []api.Field
	CreateOptions string
	TypeMap       map[reflect.Type]string
	BindPrefix    string
	Parallel      int
	MaxNameLength int
}

// newStopOrder ...
func newStopOrder() *message {
	return &message{
		Sender: "supervisor",
		Type:   stopOrder,
	}
}

// GetTempTableName ...
func (o *Options) GetTempTableName() string {

	var name string
	if o.MaxNameLength != 0 && len(o.Table) < o.MaxNameLength {
		postfixLength := o.MaxNameLength - len(o.Table)
		postFix := util.RandString(postfixLength)
		name = o.Table + postFix
	} else {
		name = util.RandString(30)
	}
	if o.TempSchema != "" {
		return fmt.Sprintf("%s.%s", o.TempSchema, name)
	}
	return name
}

// GetTableName ...
func (o *Options) GetTableName(args ...interface{}) string {
	if len(args) > 0 && args[0].(bool) {
		return o.GetTempTableName()
	}
	if o.Schema != "" {
		return fmt.Sprintf("%s.%s", o.Schema, o.Table)
	}
	return o.Table
}

// BuildCreateQuery ...
// todo check create table as select
// true, tempname
func (o *Options) BuildCreateQuery(name string) (string, error) {

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
func (o *Options) BuildInsertQuery(name string) (string, error) {

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
		name, strings.Join(columns, ","), strings.Join(values, ","),
	), nil
}

// BuildDropQuery ...
func (o *Options) BuildDropQuery(name string) (string, error) {
	return fmt.Sprintf("drop table %s", name), nil
}

// BuildTruncateQuery ...
func (o *Options) BuildTruncateQuery(name string) (string, error) {
	return fmt.Sprintf("truncate table %s", name), nil
}

// BuildUnionAllQuery ...
func (h *Helper) BuildUnionAllQuery(sources []string) string {
	var target string = h.GetTableName()
	var selects []string

	for _, s := range sources {
		selects = append(selects, fmt.Sprintf("select * from %s\n", s))
	}

	return fmt.Sprintf("insert into %s \n %s", target, strings.Join(selects, "union all\n"))
}

// DropTable ...
func (h *Helper) DropTable(name string) error {
	_, err := h.Db.Exec(h.BuildDropQuery(name))
	if err != nil {
		return err
	}
	return nil
}

// DropTables ...
func (h *Helper) DropTables(names []string) error {
	for _, name := range names {
		if err := h.DropTable(name); err != nil {
			return err
		}
	}
	return nil
}

// CreateTable ...
func (h *Helper) CreateTable(name string) error {
	_, err := h.Db.Exec(h.BuildCreateQuery(name))
	if err != nil {
		return err
	}
	return nil
}

// TruncateTable ...
func (h *Helper) TruncateTable(name string) error {
	_, err := h.Db.Exec(h.BuildTruncateQuery(name))
	if err != nil {
		return err
	}
	return nil
}

// DropCreateTable ...
func (h *Helper) DropCreateTable(name string) error {
	if err := h.DropTable(name); err != nil {
		logging.Warn("Failed to drop table", err)
	}
	if err := h.CreateTable(name); err != nil {
		logging.Error("Failed to create target table")
		return err
	}
	return nil
}

// PrepareTargetTable ...
func (h *Helper) PrepareTargetTable(table string) error {
	if strings.ToLower(h.Mode) == "create" {
		if err := h.DropCreateTable(table); err != nil {
			return err
		}
	} else if strings.ToLower(h.Mode) == "truncate" {
		if err := h.TruncateTable(table); err != nil {
			logging.Error("Failed to truncate target table")
			return err
		}
	}
	return nil
}

// func (h *Helper) {}

// Helper ...
type Helper struct {
	*Options
}

// NewHelper ...
func NewHelper(o *Options) *Helper {
	return &Helper{o}
}

// UnionAllTempTables ...
func (h *Helper) UnionAllTempTables(sources []string) error {
	query := h.BuildUnionAllQuery(sources)
	tx, err := h.Db.Begin()
	if err != nil {
		logging.Error("Failed to start transaction")
		return err
	}
	_, err = tx.Exec(query)
	if err != nil {
		logging.Error("Failed union tables")
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Run ...
func (h *Helper) Run(wire *api.Wire) error {

	var table string
	table = h.GetTableName()

	if err := h.PrepareTargetTable(table); err != nil {
		return err
	}

	var parallel int = 1
	if h.Parallel > 1 {
		parallel = h.Parallel
	}

	var in chan *message
	outs := make([]chan *message, parallel)
	for i := range outs {
		outs[i] = make(chan *message, 10) // todo buffer size from .env

	}

	tables := make([]string, parallel)
	if parallel == 1 {
		// directly to target table
		tables = append(tables, table)
		go h.insert(wire, in, outs[0], table)
	} else {
		for i := 0; i < parallel; i++ {
			tables = append(tables, h.GetTempTableName())
			if err := h.CreateTable(tables[i]); err != nil {
				stopWorkers(outs)
				return err
			}
			go h.insert(wire, in, outs[0], tables[i])
		}
	}
	if err := supervise(in, outs); err != nil {
		return err
	}
	if parallel > 1 {
		if err := h.UnionAllTempTables(tables); err != nil {
			return err
		}
	}
	return nil
}

func stopWorkers(channels []chan *message) {
	order := newStopOrder()
	for _, ch := range channels {
		ch <- order
	}
}

func supervise(in chan *message, outs []chan *message) error {
	var successCount int = 0
	for message := range in {
		if message.Type == successMessage {
			if successCount == len(outs) {
				break
			}
		} else if message.Type == errorMessage {
			killm := newStopOrder()
			for _, o := range outs {
				o <- killm
			}
			return errors.New("One of the child sessions failed")
		}
	}
	if successCount == len(outs) {
		return nil
	}
	return errors.New("Can not get success message from all childs")
}

func (h *Helper) insert(wire *api.Wire, out chan *message, in chan *message, name string) error {

	query, err := h.BuildInsertQuery(name)
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
			_, err := tx.Exec(query, d...)
			if err != nil {
				logging.Error("Failed to insert record")
				tx.Rollback()
				return err
			}
			wire.Feed <- api.NewWriteCountFeed(1)
			parentMessage, ok := <-in
			if ok {
				if parentMessage.Type == stopOrder {
					return errors.New("Terminated by the parent")
				}
			}
		}
	}
	if err := tx.Commit(); err != nil {
		out <- &message{Sender: "worker", Type: errorMessage, Content: err}
		return err
	}
	out <- &message{Sender: "worker", Type: successMessage}
	return nil
}
