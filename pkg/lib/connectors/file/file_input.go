package file

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

func (c *FileConnector) Read(d types.Dataset, w *wire.Wire) (err error) {
	files, err := c.ListFiles(d.Config)
	if err != nil {
		log.Error().Err(err).Msg("failed to list files")
		return err
	}
	d.Config["files"] = files

	return c.FileFormat.Read(d, w)
}
