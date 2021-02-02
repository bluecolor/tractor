package input

import (
	"database/sql"
	"errors"
	"fmt"
)

// Options ...
type Options struct {
	Db           *sql.DB
	Query        string
	Table        string
	Select       string
	Where        string
	DbQueryArgs  []interface{}
	BatchSize    int
	SendMetadata bool
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
func (h *Helper) BuildSelectQuery() (string, error) {
	if h.Query != "" {
		return h.Query, nil
	}
	if h.Table == "" {
		return "", errors.New("Table name is not given")
	}
	var query string
	if h.Select != "" {
		query = fmt.Sprintf("select %s from %s", h.Select, h.Table)
	} else {
		query = fmt.Sprintf("select * from %s", h.Table)
	}
	if h.Where != "" {
		query = fmt.Sprintf("%s where %s", query, h.Where)
	}

	return query, nil
}

// Run ...
func (h *Helper) Run() error {
	query, err := h.BuildSelectQuery()
	if err != nil {
		return err
	}
}

// func (o *Options) BuildQuery(args ...interface{}) (string, error) {
// 	if o.Query != "" {
// 		return o.Query, nil
// 	}
// 	fieldCount := 0
// 	if len(args) > 0 {
// 		fieldCount = args[0].(int)
// 	} else {
// 		return "", errors.New("Dynamic field resolution not supported yet")
// 	}

// 	fields := ""

// 	for i := 1; i <= fieldCount; i++ {
// 		fields = fields + ":" + strconv.Itoa(i)
// 		if i != fieldCount {
// 			fields = fields + ","
// 		}
// 	}
// 	return "insert into " + c.Table + " values(" + fields + ")", nil
// }

// // Sender ...
// type Helper struct {
// 	*Options
// }

// func NewHelper(o *Options) *Helper {
// 	return &Helper{o}
// }

// func (h *Helper) Run() error {

// }
