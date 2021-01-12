package sqlx

import (
	"database/sql"

	"github.com/bluecolor/tractor/api/md"
	"github.com/bluecolor/tractor/api/md/mdt"
	"github.com/bluecolor/tractor/api/message"
	"github.com/bluecolor/tractor/api/message/mt"
	"github.com/bluecolor/tractor/api/message/sender"
)

// SendQueryResult ...
func SendQueryResult(channel chan message.Message, db *sql.DB, query string, batchSize int, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	columns, _ := rows.ColumnTypes()
	ds, _ := getDataStore("demo", columns)
	sendMetadata(channel, ds)
	sendData(channel, len(ds.Fields), rows, batchSize)
	return nil
}

func getDataStore(name string, columns []*sql.ColumnType) (*md.DataStore, error) {

	var fields []md.Field
	for _, ct := range columns {

		precision, scale, ok := ct.DecimalSize()
		decimalSize := md.DecimalSize{Precision: precision, Scale: scale, Ok: ok}

		n, ok := ct.Nullable()
		nullable := md.Nullable{Nullable: n, Ok: ok}

		ln, ok := ct.Length()
		length := md.Length{Length: ln, Ok: ok}

		field := md.Field{
			Name:        ct.Name(),
			Type:        ct.ScanType(),
			DecimalSize: decimalSize,
			Nullable:    nullable,
			Length:      length,
		}

		fields = append(fields, field)
	}
	return &md.DataStore{Name: name, Fields: fields}, nil
}

func sendMetadata(channel chan message.Message, ds *md.DataStore) {
	md := md.Metadata{
		Type:    mdt.DataStore,
		Content: ds,
	}
	message := message.Message{
		Sender:      sender.InputPlugin,
		MessageType: mt.Metadata,
		Content:     md,
	}
	channel <- message
}

func getValuePointers(c int) []interface{} {
	var ptrs = make([]interface{}, c)
	for i := range ptrs {
		var v interface{}
		ptrs[i] = &v
	}
	return ptrs
}

func sendData(channel chan message.Message, fieldCount int, rows *sql.Rows, batchSize int) {
	valuePtrs := getValuePointers(fieldCount)

	var records [][]interface{}

	send := func(records *[][]interface{}) {
		channel <- message.NewDataMessage(message.Data{Content: *records})
	}

	for rows.Next() {
		row := make([]interface{}, fieldCount)
		if err := rows.Scan(valuePtrs...); err != nil {
			panic(err)
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
}
