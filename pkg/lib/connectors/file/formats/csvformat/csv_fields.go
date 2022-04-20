package csvformat

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"

	"go.beyondstorage.io/v5/pairs"
)

func (f *CsvFormat) ReadFields(options map[string]interface{}) ([]string, error) {
	var buf bytes.Buffer
	var lines []string
	size, offset := int64(100000), int64(0)
	files := options["files"].([]string)
	delimiter := options["delimiter"].(string)
	var readBytes int64 = -1

	if readBytes != 0 {
		readBytes, err := f.storage.Read(files[0], &buf, pairs.WithOffset(offset), pairs.WithSize(size))
		offset += readBytes
		if err != nil {
			return nil, err
		}
	}
	if readBytes == 0 {
		return nil, errors.New("no data while getting fields")
	}
	lines, _ = toLinesWithRest(buf.String())
	csvReader := csv.NewReader(strings.NewReader(strings.Join(lines, "\n")))
	csvReader.Comma = []rune(delimiter)[0]
	csvReader.LazyQuotes = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows[0], nil
}
