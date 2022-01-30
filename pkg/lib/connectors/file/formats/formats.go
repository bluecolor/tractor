package formats

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type FileFormat interface {
	Read(e meta.ExtInput, w wire.Wire) (err error)
}
