package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/bluecolor/tractor/api"
	_ "github.com/godror/godror"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
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

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte) {
	config, err := getConfig(conf)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("godror", getDataSourceName(config))
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`select * from all_tables`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		data := ""
		if err := rows.Scan(&data); err != nil {
			panic(err)
		}
		println(data)
	}
	db.Close()
	wg.Done()
}
