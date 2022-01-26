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

// todo add batch size
// todo use field-mapping instead of dataset
func (m *MySQLConnector) StartWriteWorker(e meta.ExtOutput, w wire.Wire, i int) error {
	for data := range w.ReadData() {
		query, err := m.BuildBatchInsertQuery(e.Dataset, len(data))
		if err != nil {
			w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			return err
		}
		values := make([]interface{}, len(data)*len(e.Dataset.Fields))
		for i, r := range data {
			for j, f := range e.Dataset.Fields {
				values[i*len(e.Dataset.Fields)+j] = r[f.Name]
			}
		}
		_, err = m.db.Exec(query, values...)
		if err != nil {
			w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			return err
		} else {
			w.SendFeed(feeds.NewWriteProgress(len(data)))
		}
	}
	return nil
}
func (m *MySQLConnector) Write(e meta.ExtOutput, w wire.Wire) (err error) {
	var parallel int = 1
	if e.Parallel > 1 {
		log.Warn().Msgf("parallel write is not supported for MySQL connector. Using %d", parallel)
	}
	if e.Parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", e.Parallel, parallel)
	}
	if err = m.PrepareTable(e); err != nil {
		w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
		return
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := m.StartWriteWorker(e, w, i)
			if err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			}
		}(wg, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderOutputConnector))
	return
}
func (m *MySQLConnector) BuildCreateQuery(d meta.Dataset) (query string, err error) {
	columns := ""
	for _, f := range d.Fields {
		columns += f.Name + " " + f.Type + ",\n"
	}
	columns = strings.TrimSuffix(columns, ",\n")
	query = "CREATE TABLE IF NOT EXISTS " + d.Name + " (\n" + columns + "\n);"
	return
}
func (m *MySQLConnector) BuildTruncateQuery(d meta.Dataset) (query string) {
	query = "TRUNCATE TABLE " + d.Name + ";"
	return
}
func (m *MySQLConnector) BuildDropQuery(d meta.Dataset) (query string) {
	query = "DROP TABLE IF EXISTS " + d.Name + ";"
	return
}
func (m *MySQLConnector) BuildBatchInsertQuery(d meta.Dataset, recordCount int) (query string, err error) {
	if d.Fields == nil || len(d.Fields) == 0 {
		err = fmt.Errorf("no fields found for dataset %s", d.Name)
		return
	}
	if recordCount == 0 {
		err = fmt.Errorf("no records found for dataset %s", d.Name)
		return
	}
	columns := ""
	values := ""
	for _, f := range d.Fields {
		columns += f.Name + ", "
		values += "?,"
	}
	columns = strings.TrimSuffix(columns, ", ")
	values = strings.TrimSuffix(values, ", ")
	values = strings.Repeat("("+values+"),", recordCount)
	query = "INSERT INTO " + d.Name + " (" + columns + ") VALUES " + values

	return
}
func (m *MySQLConnector) CreateTable(d meta.Dataset) (err error) {
	query, err := m.BuildCreateQuery(d)
	if err != nil {
		return
	}
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) DropTable(d meta.Dataset) (err error) {
	query := m.BuildDropQuery(d)
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) TruncateTable(d meta.Dataset) (err error) {
	query := m.BuildTruncateQuery(d)
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) PrepareTable(e meta.ExtOutput) (err error) {
	switch e.ExtractionMode {
	case meta.ExtractionModeCreate:
		if err = m.DropTable(e.Dataset); err != nil {
			return
		}
	case meta.ExtractionModeInsert:
		if err = m.TruncateTable(e.Dataset); err != nil {
			return
		}
	}
	return
}
