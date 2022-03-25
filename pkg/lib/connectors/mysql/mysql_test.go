package mysql

import (
	"database/sql"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/types"
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
	dataset := types.Dataset{
		Name: "test",
		Fields: []*types.Field{
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

	m := MySQLConnector{}
	query, err := m.BuildReadQuery(dataset, 0)
	if err != nil {
		t.Error(err)
	}
	expected := "SELECT id,name FROM test"
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
}

func TestIO(t *testing.T) {
	db, mock, err := sqlmock.New()

	mock.MatchExpectationsInOrder(false)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}

	data := [][]interface{}{
		{1, "John Doe"},
		{2, "Jane Doe"},
	}
	expectedrc := len(data)

	inputDataset := types.Dataset{
		Name: "test_in",
		Fields: []*types.Field{
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
	outputDataset := types.Dataset{
		Name: "test_out",
		Fields: []*types.Field{
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

	connector := &MySQLConnector{
		db: db,
	}
	query, err := connector.BuildReadQuery(inputDataset, 0)
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, row := range data {
		rows.AddRow(row[0], row[1])
	}

	mock.ExpectExec("DROP TABLE IF EXISTS test_out").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("TRUNCATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS test_out").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(query).WillReturnRows(rows)

	query, err = connector.BuildBatchInsertQuery(outputDataset, expectedrc)
	if err != nil {
		t.Error(err)
	}
	values := strings.Repeat("(?,?),", expectedrc)
	expected := "INSERT INTO test_out (id,full_name) VALUES " + values[:len(values)-1]
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
	mock.ExpectExec("^INSERT INTO test_out").
		WithArgs(utils.TwoToOneDim(data)...).
		WillReturnResult(sqlmock.NewResult(0, 2))

	w := wire.New()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w)
		memo := casette.Memo()
		if memo.HasError() {
			for _, e := range memo.Errors() {
				t.Error(e.Content.(error))
			}
			return
		}
		if memo.ReadCount() != expectedrc {
			t.Errorf("(read) expected %d records, got %d", expectedrc, memo.ReadCount())
			return
		}
		if memo.WriteCount() != expectedrc {
			t.Errorf("(write) expected %d records, got %d", expectedrc, memo.WriteCount())
			return
		}
	}(wg, t)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p types.Dataset, w *wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, outputDataset, w)

	// start input connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Input, p types.Dataset, w *wire.Wire) {
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, inputDataset, w)

	wg.Wait()
}

func TestGetDatabases(t *testing.T) {
	db, mock, err := sqlmock.New()

	mock.MatchExpectationsInOrder(false)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}
	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("db1")
	rows.AddRow("db2")
	rows.AddRow("db3")
	mock.ExpectQuery("SHOW DATABASES").WillReturnRows(rows)
	c := &MySQLConnector{
		db: db,
	}
	dbs, err := c.getDatabases(nil)
	if err != nil {
		t.Error(err)
	}
	names := dbs.([]string)
	if len(names) != 3 {
		t.Errorf("expected 3 databases, got %d", len(names))
	}
	if names[0] != "db1" {
		t.Errorf("expected db1, got %s", names[0])
	}
	if names[1] != "db2" {
		t.Errorf("expected db2, got %s", names[1])
	}
	if names[2] != "db3" {
		t.Errorf("expected db3, got %s", names[2])
	}
}

func TestGetTables(t *testing.T) {
	db, mock, err := sqlmock.New()

	mock.MatchExpectationsInOrder(false)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mock db")
	}
	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("table1")
	rows.AddRow("table2")
	rows.AddRow("table3")
	mock.ExpectQuery("SHOW TABLES FROM db1").WillReturnRows(rows)
	c := &MySQLConnector{
		db: db,
	}
	dbs, err := c.getTables(map[string]interface{}{"database": "db1"})
	if err != nil {
		t.Error(err)
	}
	names := dbs.([]string)
	if len(names) != 3 {
		t.Errorf("expected 3 tables, got %d", len(names))
	}
}
