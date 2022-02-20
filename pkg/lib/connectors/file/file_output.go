package file

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Write(p types.SessionParams, w *wire.Wire) (err error) {
	return c.FileFormat.Write(p, w)
}
