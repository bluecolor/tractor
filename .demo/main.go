package main

import (
	"database/sql"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
)

func main() {
	db, mock, _ := sqlmock.New()

	mock.MatchExpectationsInOrder(false)
	data := []interface{}{
		1,
		"John Doe",
		2,
		"Jane Doe",
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(data[0], data[1])
	rows.AddRow(data[2], data[3])

	mock.ExpectExec("DROP TABLE IF EXISTS test_out").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("TRUNCATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS test_out").WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery("select").WillReturnRows(rows).WithArgs()

	// mock.ExpectPrepare("INSERT INTO").
	// 	ExpectExec().
	// 	WithArgs(
	// 		data[0],
	// 		data[1],
	// 		data[2],
	// 		data[3],
	// 	).
	// 	WillReturnResult(sqlmock.NewResult(0, 2))

	mock.ExpectBegin()
	mock.
		ExpectExec("INSERT INTO").
		WithArgs(
			data[0],
			data[1],
			data[2],
			data[3],
		).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	funcs := []func(db *sql.DB) error{
		runSelect,
		runCreate,
		runInsert,
		runDrop,
		runTruncate,
	}

	// println("Testing")
	// for _, f := range funcs {
	// 	if err := f(db); err != nil {
	// 		panic(err)
	// 	}
	// }
	println("Testing with go routines")
	wg := &sync.WaitGroup{}
	for _, f := range funcs {
		wg.Add(1)
		go func(f func(db *sql.DB) error, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := f(db); err != nil {
				panic(err)
			}
		}(f, wg)
	}
	wg.Wait()
}

func runSelect(db *sql.DB) error {
	_, err := db.Query("select * from test_out")
	return err
}

func runCreate(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS test_out (a varchar(255)")
	return err
}

func runInsert(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO test_out (a, b, c, d) VALUES (?, ?, ?, ?)", 1, "John Doe", 2, "Jane Doe")
	if err != nil {
		return err
	}
	return tx.Commit()

	// _, err := db.Exec("INSERT INTO test_out (id,full_name) VALUES (?,?),(?,?)",
	// 	1, "John Doe", 2, "Jane Doe",
	// )

	// stmt, err := db.Prepare("INSERT INTO test_out (id,full_name) VALUES (?,?),(?,?)")
	// if err != nil {
	// 	return err
	// }
	// _, err = stmt.Exec(1, "John Doe", 2, "Jane Doe")
	// return err
}

func runDrop(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS test_out")
	return err
}

func runTruncate(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE test_out")
	return err
}
