package csvformat

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/formats"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"go.beyondstorage.io/v5/types"
)

type csvconfig struct {
	Delimiter string `json:"delimiter"`
	Header    bool   `json:"header"`
	Quoted    bool   `json:"quoted"`
}

type options struct {
	csvconfig  csvconfig
	bufferSize int
	fields     []meta.Field
	parallel   int
}

type CsvFormat struct {
	storage types.Storager
}

func New(storage types.Storager) (*CsvFormat, error) {
	return &CsvFormat{storage: storage}, nil
}

func init() {
	formats.Add("csv", func(storage types.Storager) (formats.FileFormat, error) {
		return New(storage)
	})
}
