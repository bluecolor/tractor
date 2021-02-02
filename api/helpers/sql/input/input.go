package input

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/api"
)

// Options ...
type Options struct {
	Db           *sql.DB
	QueryArgs    []interface{}
	BatchSize    int
	SendMetadata bool
	Queries      []string
}

// Helper ...
type Helper struct {
	*Options
}

// NewHelper ...
func NewHelper(o *Options) *Helper {
	return &Helper{o}
}

// BuildSelectQuery ...
// table, select, where
func BuildSelectQuery(args ...string) (string, error) {

	var query, table, sel, where string
	sel = "select *"
	if len(args) == 0 {
		return "", errors.New("Missing parameters in select query builder")
	}
	table = args[0]
	if len(args) > 1 {
		sel = fmt.Sprintf("select %", args[1])
	}
	if len(args) > 2 {
		where = fmt.Sprintf("where %s", args[2])
	}
	query = strings.Trim(fmt.Sprintf("%s from % %s", sel, table, where))

	return query, nil
}

// Run ...
func (h *Helper) Run(wire *api.Wire) error {

	for _, query := range h.Queries {
		go h.run(wire, query)
	}
}

func (h *Helper) run(wire *api.Wire) error {

}
