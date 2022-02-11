package file

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Read(p params.ExtParams, w *wire.Wire) (err error) {
	return c.FileFormat.Read(p, w)
}
