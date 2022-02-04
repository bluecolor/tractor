package repo

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func genSQLMockAnyArg(count int) []driver.Value {
	args := make([]driver.Value, count)
	for i := 0; i < count; i++ {
		args[i] = sqlmock.AnyArg()
	}
	return args
}

func prepareMock(mock sqlmock.Sqlmock) {
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).
			AddRow("5.7.25-0ubuntu0.18.04.1"))
}

func getRepository(db *sql.DB) (*Repository, error) {
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
	return &Repository{DB: gdb}, nil
}
func getMockRepo() (*Repository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	prepareMock(mock)
	r, err := getRepository(db)
	if err != nil {
		return nil, nil, err
	}
	return r, mock, nil
}

func TestDrop(t *testing.T) {
	repo, mock, err := getMockRepo()
	if err != nil {
		t.Fatal(err)
	}
	stmt := &gorm.Statement{DB: repo.DB}
	for _, model := range models.Models {
		stmt.Parse(model)
		tablename := stmt.Schema.Table
		mock.ExpectQuery("SELECT(.+?)FROM `" + tablename + "`")
		mock.ExpectBegin()
		mock.ExpectPrepare("DROP TABLE `" + tablename + "`").
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

	}
	repo.Drop()
}
func TestMigrate(t *testing.T) {
	repo, mock, err := getMockRepo()
	if err != nil {
		t.Fatal(err)
	}
	stmt := &gorm.Statement{DB: repo.DB}
	for _, model := range models.Models {
		stmt.Parse(model)
		tablename := stmt.Schema.Table
		mock.ExpectBegin()
		mock.ExpectPrepare("CREATE TABLE `" + tablename + "`").
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

	}
	repo.Drop()
}

// todo
func TestSeed(t *testing.T) {
}
