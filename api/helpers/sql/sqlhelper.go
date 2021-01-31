package sqlhelper

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/logging"
)

// SendOptions ...
type SendOptions struct {
	Db           *sql.DB
	Query        string
	DbQueryArgs  []interface{}
	BatchSize    int
	SendMetadata bool
}

// ReceiveOptions ...
type ReceiveOptions struct {
	Db         *sql.DB
	Table      string
	Query      string
	BindPrefix string
}

// BuildQuery ...
func (o *ReceiveOptions) BuildQuery(fieldCount int) string {
	if o.Query != "" {
		return o.Query
	}
	if len(args) > 0 {
		fieldCount = args[0].(int)
	} else {
		return "", errors.New("Dynamic field resolution not supported yet")
	}

	fields := ""

	for i := 1; i <= fieldCount; i++ {
		fields = fields + ":" + strconv.Itoa(i)
		if i != fieldCount {
			fields = fields + ","
		}
	}
	return "insert into " + c.Table + " values(" + fields + ")", nil
}

// ExecuteQuery ...
func (o *SendOptions) ExecuteQuery() (*sql.Rows, error) {
	return o.Db.Query(o.Query, o.DbQueryArgs...)
}

// Validate ...
// todo
func (o *SendOptions) Validate() error {
	return nil
}

// Send ...
func Send(wire *api.Wire, options *SendOptions) error {
	rows, err := options.ExecuteQuery()
	if err != nil {
		return err
	}

	md, err := GetFieldsMetadata(rows)
	if err != nil {
		return err
	}
	if options.SendMetadata {
		SendMetadata(wire, md)
	}

	return SendData(wire, len(md.Content.([]api.Field)), rows, options.BatchSize)
}

// GetFieldsMetadata ...
func GetFieldsMetadata(rows *sql.Rows) (*api.Metadata, error) {
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var fields []api.Field
	for _, c := range columns {

		precision, scale, ok := c.DecimalSize()
		decimalSize := api.DecimalSize{Precision: precision, Scale: scale, Ok: ok}

		n, ok := c.Nullable()
		nullable := api.Nullable{Nullable: n, Ok: ok}

		ln, ok := c.Length()
		length := api.Length{Length: ln, Ok: ok}

		field := api.Field{
			Name:        c.Name(),
			Type:        c.ScanType(),
			DecimalSize: decimalSize,
			Nullable:    nullable,
			Length:      length,
		}

		fields = append(fields, field)
	}

	md := api.Metadata{
		Type:    api.FieldsMetadata,
		Content: fields,
	}

	return &md, nil
}

// SendMetadata ...
func SendMetadata(wire *api.Wire, md *api.Metadata) {
	wire.Metadata <- md
}

// SendData ...
func SendData(wire *api.Wire, fieldCount int, rows *sql.Rows, batchSize int) error {

	var valuePtrs = make([]interface{}, fieldCount)
	for i := range valuePtrs {
		var v interface{}
		valuePtrs[i] = &v
	}

	var records [][]interface{}

	send := func(records *[][]interface{}) {
		d := api.Data{Content: *records}
		wire.Data <- &d
	}

	for rows.Next() {
		wire.Feed <- api.NewReadCountFeed(1)
		row := make([]interface{}, fieldCount)
		if err := rows.Scan(valuePtrs...); err != nil {
			logging.Error(err)
			return err
		}
		for i := 0; i < fieldCount; i++ {
			row[i] = *(valuePtrs[i].(*interface{}))
		}

		records = append(records, row)
		if len(records) == batchSize {
			send(&records)
			records = nil
		}
	}
	if len(records) > 0 {
		send(&records)
		records = nil
	}
	return nil
}

// Truncate ...
func Truncate(db *sql.DB, table string) error {
	query := fmt.Sprintf("truncate table %s", table)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Receive ...
func Receive(wire *api.Wire, options *ReceiveOptions) {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return err
	}

	var query string

	isOpen := struct {
		MetadataChannel bool
		DataChannel     bool
		FeedChannel     bool
	}{MetadataChannel: true, DataChannel: true, FeedChannel: true}

	for {
		select {
		case md, ok := <-wire.Metadata:
			if !ok {
				isOpen.MetadataChannel = false
			} else if md.Type == api.FieldsMetadata {
				query, err = cfg.BuildQuery(len(md.Content.([]api.Field)))
				if err != nil {
					return err
				}
			}
		case data, ok := <-wire.Data:
			if !ok {
				isOpen.DataChannel = false
			} else {
				if query == "" {
					query, err = cfg.BuildQuery(len(data.Content[0]))
					if err != nil {
						return nil
					}
				}
				for _, d := range data.Content {
					_, err = tx.Exec(query, d...)
					if err != nil {
						logging.Error(err)
						tx.Rollback()
						return err
					}
					wire.Feed <- api.NewWriteCountFeed(1)
				}
			}
		}
		if !isOpen.DataChannel {
			break
		}
	}

	tx.Commit()
}
