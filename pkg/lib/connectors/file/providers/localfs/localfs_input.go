package localfs

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func (f *LocalFSProvider) Read(format string, e meta.ExtInput, w wire.Wire) (err error) {

	files, err := f.FindFiles(e.Dataset.Name)
	if err != nil || len(files) == 0 {
		return err
	}
	e.Dataset.Config["files"] = files

	switch format {
	case "csv":
		err = f.ReadCsv(e, w)
	}
	return
}
