package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/message"
	"github.com/bluecolor/tractor/api/message/mt"
	"github.com/bluecolor/tractor/api/message/sender"
	"github.com/bluecolor/tractor/api/schema"
	"github.com/godror/godror"
	_ "github.com/godror/godror"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
	FetchSize        int    `yaml:"fetch_size"`
}

func (c *config) GetFetchSize() int {
	if c.FetchSize > 0 {
		return c.FetchSize
	}
	return 1000
}

// SampleConfig ...
func SampleConfig() string {
	return `
plugin: oracle
username: username
password: password
libdir: path/to/oracle_instantclient
connection_string: "localhost:1521/orcl"
table: my_table
schema:
    name: datastore_name
    fileds:
        - name: filed_name
          type: data_type
          precision: precision
          scale: scale
`
}

// Description ...
func Description() string {
	return "Read data from oracle database"
}

// PluginType ...
func PluginType() api.PluginType {
	return api.InputPlugin
}

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, channel chan message.Message) {
	config, err := getConfig(conf)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("godror", getDataSourceName(config))
	defer db.Close()

	if err != nil {
		panic(err)
	}
	rows, err := db.Query(config.getQuery(), godror.FetchArraySize(config.GetFetchSize()))

	if err != nil {
		panic(err)
	}

	columns, _ := rows.ColumnTypes()
	ds, _ := getDataStore("demo", columns)

	sendSchemaMessage(channel, ds)
	sendData(channel, len(ds.Fields), rows, config.GetFetchSize())

	close(channel)
	wg.Done()
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

	var records [][]byte

	send := func(records *[][]byte) {
		content, _ := json.Marshal(*records)
		message := message.Message{
			Sender:      sender.InputPlugin,
			MessageType: mt.Data,
			Content:     content,
		}
		channel <- message
	}

	for rows.Next() {
		row := make([]interface{}, fieldCount)
		if err := rows.Scan(valuePtrs...); err != nil {
			panic(err)
		}
		for i := 0; i < fieldCount; i++ {
			row[i] = *(valuePtrs[i].(*interface{}))
		}

		data, _ := json.Marshal(row)
		records = append(records, data)
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

func sendSchemaMessage(channel chan message.Message, ds *schema.DataStore) {
	content, _ := json.Marshal(ds)
	message := message.Message{
		Sender:      sender.InputPlugin,
		MessageType: mt.Schema,
		Content:     content,
	}
	channel <- message
}

func getDataStore(name string, columns []*sql.ColumnType) (*schema.DataStore, error) {

	var fields []schema.Field

	for _, ct := range columns {

		precision, scale, ok := ct.DecimalSize()
		decimalSize := schema.DecimalSize{Precision: precision, Scale: scale, Ok: ok}

		n, ok := ct.Nullable()
		nullable := schema.Nullable{Nullable: n, Ok: ok}

		ln, ok := ct.Length()
		length := schema.Length{Length: ln, Ok: ok}

		field := schema.Field{
			Name:        ct.Name(),
			Type:        ct.ScanType(),
			DecimalSize: decimalSize,
			Nullable:    nullable,
			Length:      length,
		}

		fields = append(fields, field)
	}
	return &schema.DataStore{Name: name, Fields: fields}, nil
}

func getConfig(conf []byte) (*config, error) {
	config := config{}
	if err := api.ParseConfig(conf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getDataSourceName(config *config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}

func (c *config) getQuery() string {
	return fmt.Sprintf("select * from %s", c.Table)
}
