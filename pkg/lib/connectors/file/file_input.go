package file

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (c *FileConnector) Read(d types.Dataset, w *wire.Wire) (err error) {
	files, err := c.ListFiles(d.Config)
	if err != nil {
		return err
	}
	d.Config["files"] = files

	return c.FileFormat.Read(d, w)
}
