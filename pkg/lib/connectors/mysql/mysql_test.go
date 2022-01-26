package mysql

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

const TIMEOUT = 5 * time.Second

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

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery(query).WillReturnRows(rows)
	wire := wire.NewWire()
	if err != m.Read(e, wire) {
		t.Error(err)
	}

	// todo check premature success
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Debug().Msg("waiting for feed")
		defer wg.Done()
		for {
			select {
			case feed := <-wire.FeedChannel:
				log.Debug().Msgf("got feed  %v", feed.Type)
				if feed.Type == feeds.SuccessFeed {
					wg.Done()
					return
				} else if feed.Type == feeds.ErrorFeed {
					wg.Done()
					t.Error(feed.Content)
				}
			case <-time.After(TIMEOUT):
				wg.Done()
				t.Error("timeout no success feed received")
			}
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Debug().Msg("waiting for data")
		dataReceived := 0
		for {
			log.Debug().Msgf("data received: %d", dataReceived)
			select {
			case feed := <-wire.DataChannel:
				if feed == nil {
					wg.Done()
					t.Error("no data")
				} else {
					dataReceived += len(feed)
				}
			case <-time.After(TIMEOUT):
				if dataReceived < 2 {
					wg.Done()
					t.Error("missing data before timeout expected 2, got", dataReceived)
				} else if dataReceived > 2 {
					wg.Done()
					t.Error("too much data before timeout expected 2, got", dataReceived)
				}
			}
		}
	}(wg)

	wg.Wait()
}