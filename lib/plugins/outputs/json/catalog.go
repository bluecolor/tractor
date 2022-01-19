package json

import (
	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/feed"
	"github.com/bluecolor/tractor/lib/utils"
)

func (j *Json) mergeSourceCatalog(sourceCatalog *config.Catalog) {
	if sourceCatalog == nil {
		return
	}
	if j.Catalog == nil {
		j.Catalog = sourceCatalog
	} else {
		if j.Catalog.Name == "" && sourceCatalog.Name != "" {
			j.Catalog.Name = sourceCatalog.Name
		}
		if sourceCatalog.Fields == nil {
			return
		}
		if j.Catalog.Fields == nil {
			j.Catalog.Fields = sourceCatalog.Fields
			return
		}
		if !j.Catalog.AutoMapFields {
			return
		}
		fieldMap := j.Catalog.GetFieldMap()
		for _, field := range sourceCatalog.Fields {
			if _, ok := fieldMap[field.Name]; !ok {
				j.Catalog.Fields = append(j.Catalog.Fields, field)
			} else {
				utils.Merge(&field, fieldMap[field.Name])
			}
		}
	}
}

func (j *Json) mapRecord(record *feed.Record) (err error) {
	return utils.MapRecord(j.Catalog, record)
}
