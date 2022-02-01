package file

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (f *FileConnector) Write(p meta.ExtParams, w wire.Wire) (err error) {
	return f.FileFormat.Write(p, w)
}
