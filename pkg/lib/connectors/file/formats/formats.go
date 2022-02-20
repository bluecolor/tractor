package formats

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type FileFormat interface {
	Read(e types.SessionParams, w *wire.Wire) (err error)
	Write(e types.SessionParams, w *wire.Wire) (err error)
}
