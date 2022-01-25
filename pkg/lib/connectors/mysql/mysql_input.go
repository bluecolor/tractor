package mysql

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feed"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

func (m *MySQLConnector) StartReadWorker(e meta.ExtInput, w wire.Wire, i int) (err error) {
	query, err := m.BuildReadQuery(e, i)
	if err != nil {
		return err
	}
	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	bufferSize := e.Config.GetInt("buffer_size", 100)
	buffer := []feed.Record{}
	for rows.Next() {
		columns := make([]interface{}, len(e.Fields))
		columnPointers := make([]interface{}, len(e.Fields))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}
		record := feed.Record{}
		for i, f := range e.Fields {
			if f.Name == "" {
				f.Name = fmt.Sprintf("col%d", i)
			}
			record[f.Name] = columns[i]
		}
		if len(buffer) >= bufferSize {
			w.SendData(buffer)
			w.SendFeed(feed.NewReadProgress(len(buffer)))
			buffer = []feed.Record{}
		} else {
			buffer = append(buffer, record)
		}
	}
	if len(buffer) > 0 {
		w.SendData(buffer)
		w.SendFeed(feed.NewReadProgress(len(buffer)))
	}
	return
}
func (m *MySQLConnector) Read(e meta.ExtInput, w wire.Wire) (err error) {
	var parallel int = 1
	if e.Parallel > 1 {
		log.Warn().Msgf("parallel read is not supported for MySQL connector. Using %d", parallel)
	}
	if e.Parallel < 1 {
		log.Warn().Msgf("invalid parallel read setting %d. Using %d", e.Parallel, parallel)
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := m.StartReadWorker(e, w, i)
			if err != nil {
				w.SendFeed(feed.NewErrorFeed(feed.SenderInputConnector, err))
			}
		}(wg, i)
	}
	return
}
func (m *MySQLConnector) BuildReadQuery(e meta.ExtInput, i int) (query string, err error) {
	if e.Fields == nil || len(e.Fields) == 0 {
		return "", fmt.Errorf("no fields specified")
	}
	columns := ""
	for _, f := range e.Fields {
		columns += f.Name + ",\n"
	}
	columns = strings.TrimRight(columns, ",\n")
	query = fmt.Sprintf("select\n%s\nfrom %s", columns, e.Dataset.Name)
	log.Debug().Msgf("query: %s", query)
	return
}
