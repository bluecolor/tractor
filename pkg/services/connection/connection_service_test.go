package connection

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/bluecolor/tractor/pkg/test"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getRepository(db *sql.DB) (*repo.Repository, error) {
	dialect := mysql.New(mysql.Config{
		DSN:        "sqlmock_db_0",
		DriverName: "mysql",
		Conn:       db,
	})
	lg := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	gdb, err := gorm.Open(dialect, &gorm.Config{
		Logger: lg,
	})
	if err != nil {
		return nil, err
	}
	return &repo.Repository{DB: gdb}, nil
}
func TestFindConnectionTypes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery("^SELECT(.+?)FROM `connection_types`").WillReturnRows(rows)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FindConnectionTypes)
	req, err := http.NewRequest(http.MethodGet, "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := []models.ConnectionType{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("handler returned unexpected body: got %v want %v", len(result), 2)
	}
}
func TestCreateConnection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).
			AddRow("5.7.25-0ubuntu0.18.04.1"))

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO").ExpectExec().
		WithArgs(test.GenSQLMockAnyArg(8)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.CreateConnection)

	var b bytes.Buffer
	connection := models.Connection{}
	gofakeit.Struct(&connection)
	if err = json.NewEncoder(&b).Encode(connection); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhsot", bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Connection{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != connection.ID {
		t.Errorf("handler returned unexpected body: got conn id %v want %v", result.ID, connection.ID)
	}
}
func TestOneConnection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `connections` WHERE").WillReturnRows(rows)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.OneConnection)
	req, err := http.NewRequest(http.MethodGet, "http://localhsot/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Connection{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}
func TestFindConnections(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery("^SELECT(.+?)FROM `connections`").WillReturnRows(rows)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FindConnections)
	req, err := http.NewRequest(http.MethodGet, "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := []models.ConnectionType{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("handler returned unexpected body: got %v want %v", len(result), 2)
	}
}
func TestTestConnection(t *testing.T) {
	service := NewService(&repo.Repository{})

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.TestConnection)

	var b bytes.Buffer
	conn := models.Connection{
		ConnectionType: &models.ConnectionType{Code: "dummy"},
	}
	if err := json.NewEncoder(&b).Encode(conn); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhsot", bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteConnection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `connections` WHERE").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("^DELETE FROM `connections` WHERE").ExpectExec().
		WithArgs(test.GenSQLMockAnyArg(1)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.DeleteConnection)

	req, err := http.NewRequest(http.MethodDelete, "http://localhsot/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Connection{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}

func TestCreateProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO `providers`").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	var b bytes.Buffer
	provider := models.Provider{
		Name: "name",
	}
	if err = json.NewEncoder(&b).Encode(provider); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhsot", bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.CreateProvider)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Provider{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Name != provider.Name {
		t.Errorf("handler returned unexpected body: got conn id %v want %v", result.Name, provider.Name)
	}
}
func TestFindProviders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "name 1").
		AddRow(2, "name 2")
	mock.ExpectQuery("^SELECT(.+?)FROM `providers`").WillReturnRows(rows)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FindProviders)
	req, err := http.NewRequest(http.MethodGet, "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := []models.Provider{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("handler returned unexpected body: got %v want %v", len(result), 2)
	}
}
func TestOneProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `providers`").WillReturnRows(rows)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.OneProvider)
	req, err := http.NewRequest(http.MethodGet, "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Provider{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}
func TestDeleteProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `providers` WHERE").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("^DELETE FROM `providers` WHERE").ExpectExec().
		WithArgs(test.GenSQLMockAnyArg(1)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.DeleteProvider)

	req, err := http.NewRequest(http.MethodDelete, "http://localhsot/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if json.Valid(rr.Body.Bytes()) == false {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}
	result := models.Provider{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}
func TestFindDatasets(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name", "connection_type_id", "config"}).AddRow(1, "name 1", 1, "{}")
	mock.ExpectQuery("^SELECT(.+?)FROM `connections`").WillReturnRows(rows)
	mock.ExpectQuery("^SELECT(.+?)FROM `connection_types`").WillReturnRows(sqlmock.NewRows(
		[]string{"id", "code"}).AddRow(1, "dummy"),
	)

	repository, err := getRepository(db)
	if err != nil {
		t.Error(err)
	}
	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FindDatasets)

	req, err := http.NewRequest(http.MethodDelete, "http://localhsot/{connectionID}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"connectionID"},
			Values: []string{"1"},
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
