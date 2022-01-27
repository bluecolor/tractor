package mysql

import (
	"database/sql"
	"database/sql/driver"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
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
	}
	e := meta.ExtInput{Dataset: dataset, Parallel: 1}
	m := MySQLConnector{}
	query, err := m.BuildReadQuery(e, 0)
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
	}
	m := MySQLConnector{
		db: db,
	}
	e := meta.ExtInput{Dataset: dataset, Parallel: 1}
	query, err := m.BuildReadQuery(e, 0)
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery(query).WillReturnRows(rows)
	w := wire.NewWire()

	go func(e meta.ExtInput, w wire.Wire) {
		if err := m.Read(e, w); err != nil {
			t.Error(err)
		}
	}(e, w)

	// todo check premature success
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		log.Debug().Msg("waiting for feed")
		defer wg.Done()
		for {
			select {
			case <-w.IsReadDone():
				return
			case feed := <-w.FeedChannel:
				log.Debug().Msgf("got feed  %v", feed.Type)
				if feed.Type == feeds.SuccessFeed && feed.Sender == feeds.SenderInputConnector {
					return
				} else if feed.Type == feeds.ErrorFeed {
					t.Error(feed.Content)
				}
			case <-time.After(TIMEOUT):
				t.Error("timeout no success feed received")
			}
		}
	}(wg, w)

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
					t.Error("no data")
				} else {
					dataReceived += len(feed)
				}
			case <-w.IsReadDone():
				return
			case <-time.After(TIMEOUT):
				if dataReceived < 2 {
					t.Error("missing data before timeout expected 2, got", dataReceived)
				} else if dataReceived > 2 {
					t.Error("too much data before timeout expected 2, got", dataReceived)
				}
			}
		}
	}(wg, w)

	log.Debug().Msg("waiting for done")

	wg.Wait()
}

func toOneDim(data [][]interface{}) []driver.Value {
	var result []driver.Value
	for _, row := range data {
		for _, col := range row {
			result = append(result, col)
		}
	}
	return result
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
	}
	ei := meta.ExtInput{Dataset: inputDataset, Parallel: 1}
	query, err := m.BuildReadQuery(ei, 0)

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
	}

	eo := meta.ExtOutput{Dataset: outputDataset, Parallel: 1, FieldMappings: []meta.FieldMapping{
		{
			SourceField: meta.Field{Name: "id", Type: "int"},
			TargetField: meta.Field{Name: "id", Type: "int"},
		},
		{
			SourceField: meta.Field{Name: "name", Type: "string"},
			TargetField: meta.Field{Name: "full_name", Type: "string"},
		},
	}}

	query, err = m.BuildBatchInsertQuery(eo.Dataset, 2)
	if err != nil {
		t.Error(err)
	}
	expected := "INSERT INTO test_out (id,full_name) VALUES (?,?),(?,?)"
	if query != expected {
		t.Errorf("query is not correct: expected %s, got %s", expected, query)
	}
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(toOneDim(data)...).WillReturnResult(sqlmock.NewResult(0, 2))

	w := wire.NewWire()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		if err != m.Read(ei, w) {
			t.Error(err)
		}
		log.Debug().Msg("read finished")
	}(wg, w)

	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		if err != m.Write(eo, w) {
			t.Error(err)
		}
		log.Debug().Msg("write finished")
	}(wg, w)

	wg.Wait()
}
