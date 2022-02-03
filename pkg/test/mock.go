package test

import (
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
)

func GenSQLMockAnyArg(count int) []driver.Value {
	args := make([]driver.Value, count)
	for i := 0; i < count; i++ {
		args[i] = sqlmock.AnyArg()
	}
	return args
}

func PrepareMock(mock sqlmock.Sqlmock) {
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).
			AddRow("5.7.25-0ubuntu0.18.04.1"))
}
