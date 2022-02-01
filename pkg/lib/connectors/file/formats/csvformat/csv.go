package csvformat

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/formats"
	"go.beyondstorage.io/v5/types"
)

const (
	DelimiterKey = "delimiter"
	FilesKey     = "files"
	QuotesKey    = "quotes"
	HeaderKey    = "header"
	FileNameKey  = "file_name"
)

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
