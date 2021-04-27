package csv

import (
	"strings"

	"github.com/bluecolor/tractor"
)

func (c *Csv) dataToString(data tractor.Data) string {
	buff := make([]string, len(data))
	for i, record := range data {
		rec := make([]string, len(record))
		for k, r := range record {
			var ok bool
			if rec[k], ok = r.(string); !ok {
				rec[k] = ""
			}
		}
		buff[i] = strings.Join(rec, c.ColumnDelim)
	}
	return strings.Join(buff, c.RecordDelim)
}
