package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/md"
	"github.com/bluecolor/tractor/api/message"
	"github.com/bluecolor/tractor/api/message/mt"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
}

func getConfig(conf []byte) (*config, error) {
	config := config{}
	if err := api.ParseConfig(conf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, channel chan message.Message) {
	config, err := getConfig(conf)
	if err != nil {
		panic(err)
	}

	m := <-channel
	ds, _ := getDataStore(&m)

	query := getQuery(len(ds.Fields), config.Table)
	db, _ := sql.Open("godror", getDataSourceName(config))
	stmt, _ := db.Prepare(query)

	println(query)
	var recievedCount int = 0
	for m := range channel {
		recievedCount = recievedCount + 1
		data := m.Content.(message.Data)
		for _, d := range data.Content {
			_, err = stmt.Exec(d...)
			if err != nil {
				panic(err)
			}
		}
	}

	db.Close()
	wg.Done()
}

func getDataStore(message *message.Message) (*md.DataStore, error) {
	if message.MessageType != mt.Metadata {
		return nil, errors.New("First message should be metadata")
	}
	ds := message.Content.(md.Metadata).Content.(*md.DataStore)
	return ds, nil
}

func getQuery(fieldCount int, table string) string {
	query := ""

	for i := 1; i <= fieldCount; i++ {
		query = query + ":" + strconv.Itoa(i)
		if i != fieldCount {
			query = query + ","
		}
	}
	return "insert into " + table + " values(" + query + ")"
}

func getDataSourceName(config *config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}
