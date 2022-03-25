package file

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Read(d types.Dataset, w *wire.Wire) (err error) {
	return c.FileFormat.Read(d, w)
}
