package csv

import (
	"strings"

	"github.com/bluecolor/tractor/lib/feed"
)

func (c *Csv) dataToString(data feed.Data) string {
	buff := make([]string, len(data))
	for i, record := range data {
		rec := make([]string, len(record))
		j := 0
		for _, v := range record {
			var ok bool
			if rec[j], ok = v.(string); !ok {
				rec[j] = ""
			}
			j++
		}
		buff[i] = strings.Join(rec, c.ColumnDelim)
	}
	return strings.Join(buff, c.RecordDelim)
}
