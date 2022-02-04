package extraction

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/test"
	"github.com/go-chi/chi"
)

func TestOneExtraction(t *testing.T) {
	repository, mock, err := test.GetMockRepo()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `extractions` WHERE").WillReturnRows(rows)

	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.OneExtraction)
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
	result := models.Extraction{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 && result.Name != "name 1" {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}
func TestFindExtractions(t *testing.T) {
	repository, mock, err := test.GetMockRepo()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1").AddRow(2, "name 2")
	mock.ExpectQuery("^SELECT(.+?)FROM `extractions`").WillReturnRows(rows)

	service := NewService(repository)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))
	defer ts.Close()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FindExtractions)
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
	result := []models.Extraction{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("handler returned unexpected body: got %v want %v", len(result), 2)
	}
}
func TestDeleteExtraction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	test.PrepareMock(mock)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name 1")
	mock.ExpectQuery("^SELECT(.+?)FROM `extractions` WHERE").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("^DELETE FROM `extractions` WHERE").ExpectExec().
		WithArgs(test.GenSQLMockAnyArg(1)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository, err := test.GetRepository(db)
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
	handler := http.HandlerFunc(service.DeleteExtraction)

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
	result := models.Extraction{}
	if err = json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", result.ID, 1)
	}
}
