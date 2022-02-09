package csvformat

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/utils"
)

func getFileChunks(f interface{}, p int) ([][]string, error) {
	if f == nil {
		return nil, errors.New("files not given")
	}
	files, ok := f.([]string)
	if !ok {
		return nil, errors.New("wrong type of 'files' option")
	}
	return utils.ToChunksStr(files, p), nil
}
func toLinesWithRest(bs string) ([]string, []byte) {
	lines := strings.Split(bs, "\n")
	var rest []byte = nil
	if !strings.HasSuffix(bs, "\n") {
		if len(lines) > 1 {
			lines = strings.Split(bs, "\n")
			rest = []byte(lines[len(lines)-1])
			lines = lines[:len(lines)-1]
		} else {
			rest = []byte(bs)
			lines = []string{}
		}
	}
	return lines, rest
}
func toRecord(row []string, fields []meta.Field) (msg.Record, error) {
	if len(row) != len(fields) {
		return nil, errors.New("wrong number of fields in record")
	}
	record := msg.Record{}
	for i, f := range fields {
		if f.Name == "" {
			f.Name = fmt.Sprintf("col%d", i)
		}
		record[f.Name] = row[i]
	}
	return record, nil
}
