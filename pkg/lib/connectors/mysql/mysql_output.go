package mysql

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

func (c *MySQLConnector) write(d types.Dataset, i int, data msg.Data) error {
	log.Debug().Msgf("writing %d records to %s", data.Count(), d.Name)
	ok := true
	query, err := c.BuildBatchInsertQuery(d, data.Count())
	if err != nil {
		log.Error().Err(err).Msg("failed to build batch insert query")
		return err
	}
	fields := d.Fields
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Order < fields[i].Order
	})
	values := make([]interface{}, data.Count()*len(fields))
	for i, record := range data {
		for j := range record {
			values[i*len(fields)+j] = record[j]
			if !ok {
				log.Debug().Msgf("field %d not found in record %d", j, i)
			}
		}
	}
	log.Debug().Msgf("executing mysql output query: %s", query)
	_, err = c.db.Exec(query, values...)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute batch insert query")
	}
	return err
}

// todo add batch size, buffer
// todo add timeout
func (m *MySQLConnector) StartWriteWorker(ctx context.Context, d types.Dataset, w *wire.Wire, i int) error {
	log.Debug().Msgf("starting mysql write worker %d", i)
	for {
		select {
		case data, ok := <-w.ReceiveData():
			if !ok {
				log.Debug().Msgf("mysql write worker %d received no data", i)
				return nil
			}
			if err := m.write(d, i, data); err != nil {
				return err
			}
			w.SendOutputProgress(data.Count())
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// todo transactions
func (m *MySQLConnector) Write(d types.Dataset, w *wire.Wire) (err error) {
	var parallel int = d.GetParallel()
	if parallel > 1 {
		log.Warn().Msgf("parallel write is not supported for MySQL connector. Using %d", 1)
		parallel = 1
	}
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
		parallel = 1
	}
	if err = m.PrepareTable(d); err != nil {
		w.SendOutputError(err)
		return
	}
	mwg := esync.NewWaitGroup(w, types.OutputConnector)
	for i := 0; i < parallel; i++ {
		mwg.Add(1)
		go func(mwg *esync.WaitGroup, i int, wire *wire.Wire) {
			defer mwg.Done()
			if err := m.StartWriteWorker(mwg.Context(), d, w, i); err != nil {
				mwg.HandleError(err)
			}
		}(mwg, i, w)
	}
	log.Debug().Msg("waiting for mysql output workers to finish")
	return mwg.Wait()
}
func (m *MySQLConnector) BuildCreateQuery(d types.Dataset) (query string, err error) {
	columns := ""
	for _, f := range d.Fields {
		columns += f.Name + " " + string(f.Type) + ",\n"
	}
	columns = strings.TrimSuffix(columns, ",\n")
	query = "CREATE TABLE IF NOT EXISTS " + d.Name + " (\n" + columns + "\n)"
	query = utils.Dedent(query)
	return
}
func (m *MySQLConnector) BuildTruncateQuery(d types.Dataset) (query string) {
	query = "TRUNCATE TABLE " + d.Name + ";"
	return
}
func (m *MySQLConnector) BuildDropQuery(d types.Dataset) (query string) {
	query = "DROP TABLE IF EXISTS " + d.Name
	return
}
func (m *MySQLConnector) BuildBatchInsertQuery(d types.Dataset, recordCount int) (query string, err error) {
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
func (m *MySQLConnector) CreateTable(d types.Dataset) (err error) {
	query, err := m.BuildCreateQuery(d)
	log.Debug().Msgf("executing query: %s", query)
	if err != nil {
		return
	}
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) DropTable(d types.Dataset) (err error) {
	query := m.BuildDropQuery(d)
	log.Debug().Msgf("executing query: %s", query)
	_, err = m.db.Exec(query)
	return
}
func (m *MySQLConnector) TruncateTable(d types.Dataset) (err error) {
	query := m.BuildTruncateQuery(d)
	_, err = m.db.Exec(query)
	return
}
func (c *MySQLConnector) PrepareTable(d types.Dataset) (err error) {
	switch d.GetExtractionMode("append") {
	case "create":
		if err = c.DropTable(d); err != nil {
			return
		}
		if err = c.CreateTable(d); err != nil {
			return
		}
	case "replace":
		if err = c.TruncateTable(d); err != nil {
			return
		}
	}
	return
}
