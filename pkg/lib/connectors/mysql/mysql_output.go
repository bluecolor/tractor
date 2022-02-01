package mysql

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

func (m *MySQLConnector) write(p meta.ExtParams, i int, data feeds.Data) error {
	ok := true
	dataset := *p.GetOutputDataset()
	query, err := m.BuildBatchInsertQuery(dataset, len(data))
	if err != nil {
		log.Error().Err(err).Msg("failed to build batch insert query")
		return err
	}
	values := make([]interface{}, len(data)*len(dataset.Fields))
	for i, r := range data {
		for j, f := range dataset.Fields {
			values[i*len(dataset.Fields)+j], ok = r[p.GetSourceFieldNameByTargetFieldName(f.Name)]
			if !ok {
				err = fmt.Errorf("field %s not found in record %d", f.Name, i)
				log.Error().Err(err).Msg("failed to build batch data")
				return err
			}
		}
	}
	log.Debug().Msgf("executing query: %s", query)
	stmt, err := m.db.Prepare(query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare batch insert query")
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute batch insert query")
		return err
	}
	return nil
}

// todo add batch size, buffer
// todo add timeout
func (m *MySQLConnector) StartWriteWorker(p meta.ExtParams, w wire.Wire, i int) error {
	for data := range w.ReadData() {
		if data == nil {
			break
		}
		err := m.write(p, i, data)
		if err != nil {
			w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			return err
		}
		w.SendFeed(feeds.NewWriteProgress(len(data)))
	}
	w.WriteWorkerDone()
	return nil
}
func (m *MySQLConnector) Write(p meta.ExtParams, w wire.Wire) (err error) {
	var parallel int = p.GetOutputParallel()
	if parallel > 1 {
		log.Warn().Msgf("parallel write is not supported for MySQL connector. Using %d", 1)
	}
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
	}
	if err = m.PrepareTable(p); err != nil {
		w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
		return
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, wire wire.Wire) {
			defer wg.Done()
			err := m.StartWriteWorker(p, w, i)
			if err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			}
		}(wg, i, w)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderOutputConnector))
	w.WriteDone()
	return
}
func (m *MySQLConnector) BuildCreateQuery(d meta.Dataset) (query string, err error) {
	columns := ""
	for _, f := range d.Fields {
		columns += f.Name + " " + f.Type + ",\n"
	}
	columns = strings.TrimSuffix(columns, ",\n")
	query = "CREATE TABLE IF NOT EXISTS " + d.Name + " (\n" + columns + "\n)"
	query = utils.Dedent(query)
	return
}
func (m *MySQLConnector) BuildTruncateQuery(d meta.Dataset) (query string) {
	query = "TRUNCATE TABLE " + d.Name + ";"
	return
}
func (m *MySQLConnector) BuildDropQuery(d meta.Dataset) (query string) {
	query = "DROP TABLE IF EXISTS " + d.Name
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
		columns += f.Name + ","
		values += "?,"
	}
	columns = strings.TrimSuffix(columns, ",")
	values = strings.TrimSuffix(values, ",")
	values = strings.Repeat("("+values+"),", recordCount)
	values = strings.TrimSuffix(values, ",")
	query = "INSERT INTO " + d.Name + " (" + columns + ") VALUES " + values
	return
}
func (m *MySQLConnector) CreateTable(d meta.Dataset) (err error) {
	query, err := m.BuildCreateQuery(d)
	log.Debug().Msgf("executing query: %s", query)
	if err != nil {
		return
	}
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) DropTable(d meta.Dataset) (err error) {
	query := m.BuildDropQuery(d)
	log.Debug().Msgf("executing query: %s", query)
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) TruncateTable(d meta.Dataset) (err error) {
	query := m.BuildTruncateQuery(d)
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) PrepareTable(p meta.ExtParams) (err error) {
	dataset := *p.GetOutputDataset()
	switch p.GetExtractionMode() {
	case meta.ExtractionModeCreate:
		if err = m.DropTable(dataset); err != nil {
			return
		}
		if err = m.CreateTable(dataset); err != nil {
			return
		}
	case meta.ExtractionModeInsert:
		if err = m.TruncateTable(dataset); err != nil {
			return
		}
	}
	return
}
