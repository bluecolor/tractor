package test

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluecolor/tractor/pkg/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetRepository(db *sql.DB) (*repo.Repository, error) {
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
func GetMockRepo() (*repo.Repository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	PrepareMock(mock)
	r, err := GetRepository(db)
	if err != nil {
		return nil, nil, err
	}
	return r, mock, nil
}
