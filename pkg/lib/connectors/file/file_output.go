package file

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Write(d types.Dataset, w *wire.Wire) (err error) {
	return c.FileFormat.Write(d, w)
}
