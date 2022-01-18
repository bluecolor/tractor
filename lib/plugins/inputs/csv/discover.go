package csv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bluecolor/tractor/lib/config"
)

func (c *Csv) DiscoverCatalog() (catalog *config.Catalog, err error) {
	files, err := c.getFiles()
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("given path is empty")
	}
	f, err := os.Open(files[0])
	if err != nil {
		return nil, err
	}

	catalog = &config.Catalog{
		Name: strings.ReplaceAll(f.Name(), ".", "_"),
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	columns := strings.Split(scanner.Text(), c.ColumnDelim)
	catalog.Fields = make([]config.Field, len(columns))
	for i, h := range columns {
		name := fmt.Sprintf("column_%d", i)
		if c.Header {
			name = h
		}
		catalog.Fields[i] = config.Field{
			Name: name,
			Type: "string",
		}
	}
	return catalog, nil
}
