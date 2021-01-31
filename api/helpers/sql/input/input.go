package input

import (
	"database/sql"
	"errors"
	"strconv"
)

// Options ...
type Options struct {
	Db           *sql.DB
	Query        string
	DbQueryArgs  []interface{}
	BatchSize    int
	SendMetadata bool
}

func (o *Options) BuildQuery(args ...interface{}) (string, error) {
	if o.Query != "" {
		return o.Query, nil
	}
	fieldCount := 0
	if len(args) > 0 {
		fieldCount = args[0].(int)
	} else {
		return "", errors.New("Dynamic field resolution not supported yet")
	}

	fields := ""

	for i := 1; i <= fieldCount; i++ {
		fields = fields + ":" + strconv.Itoa(i)
		if i != fieldCount {
			fields = fields + ","
		}
	}
	return "insert into " + c.Table + " values(" + fields + ")", nil
}

// Sender ...
type Helper struct {
	*Options
}

func NewHelper(o *Options) *Helper {
	return &Helper{o}
}

func (h *Helper) Run() error {

}
