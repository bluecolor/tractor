package formats

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type FileFormat interface {
	Read(d types.Dataset, w *wire.Wire) (err error)
	Write(d types.Dataset, w *wire.Wire) (err error)
	ReadFields(options map[string]interface{}) ([]string, error)
}
