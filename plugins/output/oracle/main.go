package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/message"
	"github.com/bluecolor/tractor/api/metadata"
	"github.com/bluecolor/tractor/logging"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
}

func (c *config) BuildQuery(args ...interface{}) (string, error) {
	var fieldCount int
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

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, channel chan *message.Message) error {
	config := config{}
	if err := api.ParseConfig(conf, &config); err != nil {
		return err
	}
	var fieldCount, err = getFieldCount(channel, 10)
	if err != nil {
		return err
	}

	query, err := config.BuildQuery(fieldCount)
	if err != nil {
		return err
	}
	db, _ := sql.Open("godror", getDataSourceName(&config))
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		return err
	}

	for m := range channel {
		if m.Type == message.Data {
			data := m.Content.(metadata.Data)
			for _, d := range data.Content {
				_, err = tx.Exec(query, d...)
				if err != nil {
					logging.Error(err)
					tx.Rollback()
					return err
				}
			}
		}
	}
	tx.Commit()
	db.Close()
	wg.Done()
	return nil
}

func getDataSourceName(config *config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}

func getFieldCount(channel chan *message.Message, messageLimit int) (int, error) {
	messageCount := 0
	for {
		m := <-channel
		switch m.Type {
		case message.Metadata:
			md := m.Content.(*metadata.Metadata)
			if md.Type == metadata.Fields {
				return len(md.Content.([]metadata.Field)), nil
			}
		case message.Data:
			channel <- m
			data := m.Content.(metadata.Data).Content
			return len(data[0]), nil
		}
		messageCount++
		if messageCount > messageLimit {
			return 0, errors.New("Could not find filed count within limits")
		}
	}

	return 0, errors.New("Could not find filed count")
}
