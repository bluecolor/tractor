package file

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Read(p meta.ExtParams, w *wire.Wire) (err error) {
	return c.FileFormat.Read(p, w)
}
