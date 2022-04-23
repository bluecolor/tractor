package jsonformat

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/formats"
	"go.beyondstorage.io/v5/types"
)

type JsonFormat struct {
	storage types.Storager
}

func New(storage types.Storager) (*JsonFormat, error) {
	return &JsonFormat{storage: storage}, nil
}

func init() {
	formats.Add("json", func(storage types.Storager) (formats.FileFormat, error) {
		return New(storage)
	})
}
