package csv

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/bluecolor/tractor/lib/config"
)

func (c *Csv) getFiles() (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return nil, err
	}
	for _, f := range fileInfos {
		if !f.IsDir() {
			files = append(files, path.Join(c.Path, f.Name()))
		}
	}
	return files, nil
}

func (c *Csv) toRecord(r []string) map[string]interface{} {
	record := make(map[string]interface{})
	if c.Catalog != nil && len(c.Catalog.Fields) > 0 {
		for i, column := range r {
			if i < len(c.Catalog.Fields) {
				record[c.Catalog.Fields[i].Name] = column
			}
		}
	} else {
		for i, column := range r {
			record[strings.ToLower(fmt.Sprintf("column_%d", i))] = column
		}
	}
	return record
}

func (c *Csv) headerToCatalog(columns []string) {
	name := strings.ToLower(strings.Replace(c.File, ".csv", "", 1))
	if c.Catalog == nil {
		c.Catalog = &config.Catalog{
			Fields: []config.Field{},
			Name:   name,
		}
	} else if c.Catalog.Name == "" {
		c.Catalog.Name = name
	}
	if len(c.Catalog.Fields) == 0 {
		for _, column := range columns {
			col := strings.ToLower(column)
			c.Catalog.Fields = append(c.Catalog.Fields, config.Field{
				Name: col,
			})
		}
	}
}
