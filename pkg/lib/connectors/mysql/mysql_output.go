package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

func (c *MySQLConnector) write(p meta.ExtParams, i int, data []msg.Record) error {
	ok := true
	dataset := *p.GetOutputDataset()
	query, err := c.BuildBatchInsertQuery(dataset, len(data))
	if err != nil {
		log.Error().Err(err).Msg("failed to build batch insert query")
		return err
	}
	values := make([]interface{}, len(data)*len(dataset.Fields))
	for i, r := range data {
		for j, f := range dataset.Fields {
			values[i*len(dataset.Fields)+j], ok = r[p.GetSourceFieldNameByTargetFieldName(f.Name)]
			if !ok {
				log.Debug().Msgf("field %s not found in record %d", f.Name, i)
			}
		}
	}
	_, err = c.db.Exec(query, values...)
	return err
}

// todo add batch size, buffer
// todo add timeout
func (m *MySQLConnector) StartWriteWorker(ctx context.Context, p meta.ExtParams, w *wire.Wire, i int) error {

	for {
		select {
		case msg, ok := <-w.GetDataMessage():
			if !ok {
				return nil
			}
			data := msg.Data()
			if err := m.write(p, i, data); err != nil {
				return err
			}
			w.SendOutputProgress(len(data))
		case <-w.Context().Done():
			return w.Context().Err()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// todo transactions
func (m *MySQLConnector) Write(p meta.ExtParams, w *wire.Wire) (err error) {
	var parallel int = p.GetOutputParallel()
	if parallel > 1 {
		log.Warn().Msgf("parallel write is not supported for MySQL connector. Using %d", 1)
		parallel = 1
	}
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
		parallel = 1
	}
	if err = m.PrepareTable(p); err != nil {
		w.SendOutputError(err)
		return
	}
	mwg := esync.NewManagedWaitGroup(w, types.OutputConnector)
	for i := 0; i < parallel; i++ {
		mwg.Add(1)
		go func(mwg *esync.ManagedWaitGroup, i int, wire *wire.Wire) {
			defer mwg.Done()
			if err := m.StartWriteWorker(mwg.Context(), p, w, i); err != nil {
				mwg.HandleError(err)
			}
		}(mwg, i, w)
	}
	return mwg.Wait()
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
		// if err = m.CreateTable(dataset); err != nil {
		// 	return
		// }
	case meta.ExtractionModeInsert:
		if err = m.TruncateTable(dataset); err != nil {
			return
		}
	}
	return
}
