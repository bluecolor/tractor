package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

func (m *MySQLConnector) BuildReadQuery(d types.Dataset, i int) (query string, err error) {
	fields := d.Fields
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Order < fields[i].Order
	})
	if len(fields) == 0 {
		return "", fmt.Errorf("no fields specified")
	}
	columns := ""
	for _, f := range fields {
		columns += f.GetExpressionOrName() + ","
	}
	columns = strings.TrimRight(columns, ",")
	query = fmt.Sprintf("SELECT %s FROM %s", columns, d.Name)
	log.Debug().Msgf("query: %s", query)
	return
}
func (m *MySQLConnector) StartReadWorker(d types.Dataset, w *wire.Wire, i int) (err error) {
	log.Debug().Msgf("starting mysql read worker %d", i)
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	bw := wire.NewBuffered(w, d.GetBufferSize())
	query, err := m.BuildReadQuery(d, i)
	if err != nil {
		log.Error().Err(err).Msg("failed to build read query")
		return err
	}
	log.Debug().Msgf("mysql select query: %s", query)
	if m.db == nil {
		log.Error().Msg("database connection is not initialized")
		return fmt.Errorf("database connection is not initialized")
	}
	rows, err := m.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute query")
		return err
	}
	defer rows.Close()

	for rows.Next() {
		record := make(msg.Record, len(d.Fields))
		pointers := make(msg.Record, len(d.Fields))
		for i := range record {
			pointers[i] = &record[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return err
		}
		bw.Send(record)
	}
	bw.Flush()
	return
}
func (m *MySQLConnector) Read(d types.Dataset, w *wire.Wire) (err error) {
	log.Debug().Msg("starting mysql read")
	var parallel int = d.GetParallel()
	if parallel > 1 {
		log.Warn().Msgf("parallel read is not supported for MySQL connector. Using %d", 1)
		parallel = 1
	}
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel read setting %d. Using %d", parallel, 1)
		parallel = 1
	}
	wg := esync.NewWaitGroup(w, types.InputConnector)
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *esync.WaitGroup, i int) {
			defer wg.Done()
			if err := m.StartReadWorker(d, w, i); err != nil {
				log.Error().Err(err).Msg("read worker failed")
				wg.HandleError(err)
			}
		}(wg, i)
	}
	wg.Wait()
	return
}
