package mysql

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

const TIMEOUT = 3 * time.Second

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}

	return db, mock
}
func TestBuildReadQuery(t *testing.T) {
	dataset := meta.Dataset{
		Name: "test",
		Fields: []meta.Field{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "name",
				Type: "string",
			},
		},
		Config: map[string]interface{}{
			"parallel": 1,
		},
	}
	p := meta.ExtParams{}
	p.WithInputDataset(dataset)

	m := MySQLConnector{}
	query, err := m.BuildReadQuery(p, 0)
	if err != nil {
		t.Error(err)
	}
	expected := "SELECT id,name FROM test"
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
}
func TestRead(t *testing.T) {
	db, mock := NewMock()
	dataset := meta.Dataset{
		Name: "test",
		Fields: []meta.Field{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "name",
				Type: "string",
			},
		},
		Config: map[string]interface{}{
			"parallel": 1,
		},
	}
	m := MySQLConnector{
		db: db,
	}
	p := meta.ExtParams{}.WithInputDataset(dataset)
	query, err := m.BuildReadQuery(p, 0)
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery(query).WillReturnRows(rows)
	w := wire.NewWire()

	// todo check premature success
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		log.Debug().Msg("waiting for data")
		dataReceived := 0
		for {
			log.Debug().Msgf("data received: %d", dataReceived)
			select {
			case feed := <-w.DataChannel:
				if feed == nil {
					if dataReceived < 2 {
						t.Error("missing data before timeout expected 2, got", dataReceived)
					} else if dataReceived > 2 {
						t.Error("too much data before timeout expected 2, got", dataReceived)
					}
					return
				} else {
					dataReceived += len(feed)
				}
			case <-time.After(TIMEOUT):
				t.Error("timeout before read done")
			}
		}
	}(wg, w)

	go func(p meta.ExtParams, w wire.Wire) {
		if err := m.Read(p, w); err != nil {
			t.Error(err)
		}
	}(p, w)

	log.Debug().Msg("waiting for done")

	wg.Wait()
}
func TestWrite(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}

	m := MySQLConnector{
		db: db,
	}
	data := [][]interface{}{
		{1, "name 1"},
		{2, "name 2"},
	}

	inputDataset := meta.Dataset{
		Name: "test_in",
		Fields: []meta.Field{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "name",
				Type: "string",
			},
		},
		Config: map[string]interface{}{
			"parallel": 1,
		},
	}
	ip := meta.ExtParams{}.WithInputDataset(inputDataset)
	query, err := m.BuildReadQuery(ip, 0)

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, row := range data {
		rows.AddRow(row[0], row[1])
	}
	mock.ExpectQuery(query).WillReturnRows(rows)

	outputDataset := meta.Dataset{
		Name: "test_out",
		Fields: []meta.Field{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "full_name",
				Type: "string",
			},
		},
		Config: map[string]interface{}{
			"parallel": 1,
		},
	}

	op := meta.ExtParams{}.WithOutputDataset(outputDataset).WithFieldMappings(
		[]meta.FieldMapping{
			{
				SourceField: meta.Field{Name: "id", Type: "int"},
				TargetField: meta.Field{Name: "id", Type: "int"},
			},
			{
				SourceField: meta.Field{Name: "name", Type: "string"},
				TargetField: meta.Field{Name: "full_name", Type: "string"},
			},
		},
	)

	query, err = m.BuildBatchInsertQuery(*op.GetOutputDataset(), 2)
	if err != nil {
		t.Error(err)
	}
	expected := "INSERT INTO test_out (id,full_name) VALUES (?,?),(?,?)"
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(utils.TwoToOneDim(data)...).WillReturnResult(sqlmock.NewResult(0, 2))

	w := wire.NewWire()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, ip meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err != m.Read(ip, w) {
			t.Error(err)
		}
		log.Debug().Msg("read finished")
	}(wg, ip, w)

	wg.Add(1)
	go func(wg *sync.WaitGroup, op meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err != m.Write(op, w) {
			t.Error(err)
		}
		log.Debug().Msg("write finished")
	}(wg, op, w)

	wg.Wait()
}
