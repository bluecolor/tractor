package db

import (
	"database/sql"
	"fmt"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/utils"
)

func GetFields(query string, db *sql.DB) (fields []utils.Field, err error) {
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	cts, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for _, ct := range cts {
		f := utils.Field{
			Name: ct.Name(),
		}
		if precision, scale, ok := ct.DecimalSize(); ok {
			f.Precision = precision
			f.Scale = scale
		}
		if nullable, ok := ct.Nullable(); ok {
			f.Nullable = nullable
		}
		if length, ok := ct.Length(); ok {
			f.Length = length
		}
		fields = append(fields, f)
	}

	return fields, nil
}

func GetCount(query string, db *sql.DB) (count int, err error) {
	var q string = fmt.Sprintf("select count(1) from (%s)", query)
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		return count, err
	}
	rows.Next()
	if err := rows.Scan(&count); err != nil {
		return count, err
	}
	return count, nil
}

func Read(wire tractor.Wire, query string, db *sql.DB) (err error) {
	defer func() {
		if err != nil {
			wire.SendMessage(tractor.NewErrorFeed(tractor.Anonymous, err))
		}
	}()
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	colCount := len(columns)

	var ptrs = make([]interface{}, colCount)
	for i := range ptrs {
		var v interface{}
		ptrs[i] = &v
	}
	var data tractor.Data
	send := func(data tractor.Data) {
		message := tractor.NewDataMessage(data)
		wire.SendMessage(message)
	}
	for rows.Next() {
		record := make(tractor.Record, colCount)
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		for i := 0; i < colCount; i++ {
			record[i] = *(ptrs[i].(*interface{}))
		}

		data = append(data, record)
		if len(data) >= 100 { // todo
			send(data)
			data = nil
		}
		// todo
		// order, ok := <-in
		// if ok && order.IsStopOrder() {
		// 	break
		// }
	}
	if len(data) > 0 {
		send(data)
		data = nil
	}
	wire.SendMessage(tractor.NewSuccessFeed(tractor.InputPlugin))
	return err
}
