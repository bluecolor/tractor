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

func (m *MySQLConnector) BuildReadQuery(p meta.ExtParams, i int) (query string, err error) {
	fields := p.GetFMInputFields()
	if fields == nil || len(fields) == 0 {
		return "", fmt.Errorf("no fields specified")
	}
	columns := ""
	for _, f := range fields {
		columns += f.GetExpressionOrName() + ","
	}
	columns = strings.TrimRight(columns, ",")
	query = fmt.Sprintf("SELECT %s FROM %s", columns, p.GetInputDataset().Name)
	log.Debug().Msgf("query: %s", query)
	return
}
func (m *MySQLConnector) StartReadWorker(p meta.ExtParams, w wire.Wire, i int) (err error) {
	query, err := m.BuildReadQuery(p, i)
	if err != nil {
		return err
	}
	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	bufferSize := p.GetInputBufferSize()
	buffer := []feeds.Record{}
	fields := p.GetFMInputFields()
	for rows.Next() {
		columns := make([]interface{}, len(fields))
		columnPointers := make([]interface{}, len(fields))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}
		record := feeds.Record{}
		for i, f := range fields {
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
func (m *MySQLConnector) Read(p meta.ExtParams, w wire.Wire) (err error) {
	var parallel int = p.GetInputParallel()
	if parallel > 1 {
		log.Warn().Msgf("parallel read is not supported for MySQL connector. Using %d", 1)
		parallel = 1
	}
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel read setting %d. Using %d", parallel, 1)
		parallel = 1
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			if err := m.StartReadWorker(p, w, i); err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderInputConnector, err))
			}
		}(wg, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderInputConnector))
	w.ReadDone()
	return
}
