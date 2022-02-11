package formats

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type FileFormat interface {
	Read(e params.SessionParams, w *wire.Wire) (err error)
	Write(e params.SessionParams, w *wire.Wire) (err error)
}
