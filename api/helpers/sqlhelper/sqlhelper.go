package sqlhelper

import (
	"database/sql"

	"github.com/bluecolor/tractor/api/message"
	"github.com/bluecolor/tractor/api/metadata"
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
func Send(channel chan *message.Message, options *Options) error {
	rows, err := options.ExecuteQuery()
	if err != nil {
		return err
	}

	md, err := GetFieldsMetadata(rows)
	if err != nil {
		return err
	}
	if options.SendMetadata {
		SendMetadata(channel, md)
	}

	return SendData(channel, len(md.Content.([]metadata.Field)), rows, options.BatchSize)
}

// GetFieldsMetadata ...
func GetFieldsMetadata(rows *sql.Rows) (*metadata.Metadata, error) {
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var fields []metadata.Field
	for _, c := range columns {

		precision, scale, ok := c.DecimalSize()
		decimalSize := metadata.DecimalSize{Precision: precision, Scale: scale, Ok: ok}

		n, ok := c.Nullable()
		nullable := metadata.Nullable{Nullable: n, Ok: ok}

		ln, ok := c.Length()
		length := metadata.Length{Length: ln, Ok: ok}

		field := metadata.Field{
			Name:        c.Name(),
			Type:        c.ScanType(),
			DecimalSize: decimalSize,
			Nullable:    nullable,
			Length:      length,
		}

		fields = append(fields, field)
	}

	md := metadata.Metadata{
		Type:    metadata.Fields,
		Content: fields,
	}

	return &md, nil
}

// SendMetadata ...
func SendMetadata(channel chan *message.Message, md *metadata.Metadata) {

	message := message.Message{
		Type:    message.Metadata,
		Sender:  message.InputPlugin,
		Content: md,
	}
	channel <- &message
}

// SendData ...
func SendData(channel chan *message.Message, fieldCount int, rows *sql.Rows, batchSize int) error {

	var valuePtrs = make([]interface{}, fieldCount)
	for i := range valuePtrs {
		var v interface{}
		valuePtrs[i] = &v
	}

	var records [][]interface{}

	send := func(records *[][]interface{}) {
		m := message.NewDataMessage(metadata.Data{Content: *records})
		channel <- &m
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
