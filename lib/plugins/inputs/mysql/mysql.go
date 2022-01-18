package mysql

import (
	"database/sql"

	"github.com/bluecolor/tractor/lib/config"
)

type Mysql struct {
	Host      string          `yaml:"host"`
	Port      int             `yaml:"port"`
	Database  string          `yaml:"database"`
	Username  string          `yaml:"username"`
	Password  string          `yaml:"password"`
	URL       string          `yaml:"url"`
	Query     string          `yaml:"query"`
	Select    string          `yaml:"select"`
	Where     string          `yaml:"where"`
	Table     string          `yaml:"table"`
	FetchSize int             `yaml:"fetch_size"`
	Parallel  int             `yaml:"parallel"`
	Catalog   *config.Catalog `yaml:"catalog"`
	db        *sql.DB         `yaml:"-"`
}
