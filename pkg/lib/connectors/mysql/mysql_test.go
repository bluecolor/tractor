package mysql

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/test"
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
func TestReadWrite(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}

	connector := &MySQLConnector{
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
	query, err := connector.BuildReadQuery(ip, 0)
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

	dq, _ := connector.BuildCreateQuery(outputDataset)
	mock.ExpectExec(dq).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(query).WillReturnRows(rows)

	op := meta.ExtParams{}.WithOutputDataset(outputDataset).WithFieldMappings(fm)

	query, err = connector.BuildBatchInsertQuery(*op.GetOutputDataset(), 2)
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

	w, _, cancel := wire.NewWithDefaultTimeout()
	wg := &sync.WaitGroup{}
	expectedrc := len(data)

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, c context.CancelFunc, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w, c)
		memo := casette.GetMemo()
		if memo.HasError() {
			for _, e := range memo.Errors {
				t.Error(e)
			}
		}
		if memo.ReadCount != expectedrc {
			t.Errorf("(read) expected %d records, got %d", expectedrc, memo.ReadCount)
		}
		if memo.WriteCount != expectedrc {
			t.Errorf("(write) expected %d records, got %d", expectedrc, memo.WriteCount)
		}
	}(wg, cancel, t)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, op, w)

	// start input connector
	go func(wg *sync.WaitGroup, c connectors.Input, p meta.ExtParams, w wire.Wire) {
		wg.Add(1)
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, ip, w)

	wg.Wait()
}
