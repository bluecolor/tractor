package file

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Write(p params.SessionParams, w *wire.Wire) (err error) {
	return c.FileFormat.Write(p, w)
}
