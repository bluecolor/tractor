package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/helpers/sqlhelper"
	"github.com/bluecolor/tractor/logging"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
	Truncate         bool   `yaml:"truncate"`
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
func Run(wg *sync.WaitGroup, conf []byte, wire *api.Wire) error {

	cfg := config{Truncate: true}
	if err := api.ParseConfig(conf, &cfg); err != nil {
		return err
	}

	db, err := sql.Open("godror", getDataSourceName(&cfg))
	if err != nil {
		logging.Error(err)
		return err
	}
	if err := sqlhelper.Truncate(db, cfg.Table); err != nil {
		logging.Error(err)
		return err
	}

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
	db.Close()
	wg.Done()
	return nil
}

func getDataSourceName(config *config) string {
	return fmt.Sprintf(`user="%s" password="%s" connectString="%s" libDir="%s"`,
		config.Username, config.Password, config.ConnectionString, config.Libdir)
}
