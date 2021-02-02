package input

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/api"
	helper "github.com/bluecolor/tractor/api/helpers"
	"github.com/bluecolor/tractor/api/helpers/message"
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
	helper.Supervisor
}

// NewHelper ...
func NewHelper(o *Options) *Helper {
	return &Helper{o, helper.Supervisor{}}
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
		sel = fmt.Sprintf("select %s", args[1])
	}
	if len(args) > 2 {
		where = fmt.Sprintf("where %s", args[2])
	}
	query = strings.Trim(fmt.Sprintf("%s from %s %s", sel, table, where), " ")

	return query, nil
}

// Run ...
func (h *Helper) Run(wire *api.Wire) error {

	var in chan *message.Message
	outs := make([]chan *message.Message, len(h.Queries))
	for i := range outs {
		outs[i] = make(chan *message.Message, 10) // todo buffer size from .env
	}
	for i, query := range h.Queries {
		go h.run(wire, in, outs[i], query)
	}

	return h.Supervise(in, outs)
}

func (h *Helper) run(wire *api.Wire, out chan *message.Message, in chan *message.Message, query string) error {
	return nil
}
