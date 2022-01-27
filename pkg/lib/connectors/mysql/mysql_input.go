package mysql

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
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
	buffer := []feeds.Record{}
	for rows.Next() {
		columns := make([]interface{}, len(e.Fields))
		columnPointers := make([]interface{}, len(e.Fields))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}
		record := feeds.Record{}
		for i, f := range e.Fields {
			if f.Name == "" {
				f.Name = fmt.Sprintf("col%d", i)
			}
			record[f.Name] = columns[i]
		}
		if len(buffer) >= bufferSize {
			w.SendData(buffer)
			w.SendFeed(feeds.NewReadProgress(len(buffer)))
			buffer = []feeds.Record{}
		} else {
			buffer = append(buffer, record)
		}
	}
	if len(buffer) > 0 {
		w.SendData(buffer)
		w.SendFeed(feeds.NewReadProgress(len(buffer)))
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
			if err := m.StartReadWorker(e, w, i); err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderInputConnector, err))
			}
		}(wg, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderInputConnector))
	w.ReadDone()
	return
}
func (m *MySQLConnector) BuildReadQuery(e meta.ExtInput, i int) (query string, err error) {
	if e.Fields == nil || len(e.Fields) == 0 {
		return "", fmt.Errorf("no fields specified")
	}
	columns := ""
	for _, f := range e.Fields {
		columns += f.GetExpressionOrName() + ","
	}
	columns = strings.TrimRight(columns, ",")
	query = fmt.Sprintf("SELECT %s FROM %s", columns, e.Dataset.Name)
	log.Debug().Msgf("query: %s", query)
	return
}
