package db

import (
	"database/sql"
	"fmt"
	"strings"

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
			wire.SendFeed(tractor.NewErrorFeed(tractor.Anonymous, err))
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
			wire.SendData(data)
			data = nil
		}
		// todo
		// order, ok := <-in
		// if ok && order.IsStopOrder() {
		// 	break
		// }
	}
	if len(data) > 0 {
		wire.SendData(data)
		data = nil
	}
	wire.SendFeed(tractor.NewSuccessFeed(tractor.InputPlugin))
	return err
}

func DropTable(db *sql.DB, table string) error {
	_, err := db.Exec(fmt.Sprintf("drop table %s", table))
	if err != nil {
		return err
	}
	return nil
}

func CreateTable(db *sql.DB, table string, columns []string, props string) error {
	q := fmt.Sprintf(`create table %s (
        %s
    ) %s`, table, strings.Join(columns, ", "), props)
	_, err := db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func Insert(tx *sql.Tx, query string, data *tractor.Data) (count int, err error) {
	for _, record := range *data {
		_, err = tx.Exec(query, record...)
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
