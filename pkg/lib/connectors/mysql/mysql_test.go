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
	fm := []meta.FieldMapping{
		{
			SourceField: meta.Field{Name: "id"},
			TargetField: meta.Field{Name: "id"},
		},
		{
			SourceField: meta.Field{Name: "name"},
			TargetField: meta.Field{Name: "name"},
		},
	}

	p := meta.ExtParams{}.
		WithInputDataset(dataset).
		WithFieldMappings(fm)

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
	fm := []meta.FieldMapping{
		{
			SourceField: meta.Field{Name: "id"},
			TargetField: meta.Field{Name: "id"},
		},
		{
			SourceField: meta.Field{Name: "name"},
			TargetField: meta.Field{Name: "name"},
		},
	}
	m := MySQLConnector{
		db: db,
	}
	p := meta.ExtParams{}.WithInputDataset(dataset).WithFieldMappings(fm)
	query, err := m.BuildReadQuery(p, 0)
	if err != nil {
		t.Error(err)
	}
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery(query).WillReturnRows(rows)
	w := wire.New()

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
			case feed := <-w.ReadData():
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
	fm := []meta.FieldMapping{
		{
			SourceField: meta.Field{Name: "id"},
			TargetField: meta.Field{Name: "id"},
		},
		{
			SourceField: meta.Field{Name: "name"},
			TargetField: meta.Field{Name: "full_name", Type: "string"},
		},
	}
	ip := meta.ExtParams{}.WithInputDataset(inputDataset).WithFieldMappings(fm)
	query, err := m.BuildReadQuery(ip, 0)
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, row := range data {
		rows.AddRow(row[0], row[1])
	}
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

	mock.ExpectExec("DROP TABLE IF EXISTS test_out").WillReturnResult(sqlmock.NewResult(0, 0))

	dq, _ := m.BuildCreateQuery(outputDataset)
	mock.ExpectExec(dq).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(query).WillReturnRows(rows)

	op := meta.ExtParams{}.WithOutputDataset(outputDataset).WithFieldMappings(fm)

	query, err = m.BuildBatchInsertQuery(*op.GetOutputDataset(), 2)
	if err != nil {
		t.Error(err)
	}
	expected := "INSERT INTO test_out (id,full_name) VALUES (?,?),(?,?)"
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(utils.TwoToOneDim(data)...).
		WillReturnResult(sqlmock.NewResult(0, 2))

	w := wire.New()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, ip meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err := m.Read(ip, w); err != nil {
			t.Error(err)
		}
		log.Debug().Msg("read finished")
	}(wg, ip, w)

	wg.Add(1)
	go func(wg *sync.WaitGroup, op meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err := m.Write(op, w); err != nil {
			t.Error(err)
		}
		log.Debug().Msg("write finished")
	}(wg, op, w)

	wg.Wait()
}
