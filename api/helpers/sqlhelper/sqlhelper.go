package sqlhelper

import (
	"database/sql"
	"fmt"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/logging"
)

// Options ...
type Options struct {
	Db           *sql.DB
	Query        string
	DbQueryArgs  []interface{}
	BatchSize    int
	SendMetadata bool
}

// ExecuteQuery ...
func (o *Options) ExecuteQuery() (*sql.Rows, error) {
	return o.Db.Query(o.Query, o.DbQueryArgs...)
}

// Validate ...
func (o *Options) Validate() error {
	return nil
}

// Send ...
func Send(wire *api.Wire, options *Options) error {
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
		println(len(*records))
		d := api.Data{Content: records}
		wire.Data <- &d
	}

	for rows.Next() {
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
