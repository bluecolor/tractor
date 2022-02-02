package connection

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getRepository(db *sql.DB) (*repo.Repository, error) {
	dialect := mysql.New(mysql.Config{
		DSN:        "sqlmock_db_0",
		DriverName: "mysql",
		Conn:       db,
	})
	gdb, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &repo.Repository{DB: gdb}, nil
}

func TestTestConnection(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"id"}))

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
	_, err = http.NewRequest("GET", "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.TestConnection)

	req, err := http.NewRequest("GET", "http://localhsot", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateConnection(t *testing.T) {
}
