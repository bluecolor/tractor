package formats

import (
	"fmt"

	"go.beyondstorage.io/v5/types"
)

type Creator func(storage types.Storager) (FileFormat, error)

var FileFormats = map[string]Creator{}

func Add(name string, creator Creator) {
	FileFormats[name] = creator
}
func GetFileFormat(name string, storage types.Storager) (FileFormat, error) {
	if creator, ok := FileFormats[name]; ok {
		return creator(storage)
	}
	return nil, fmt.Errorf("file format %s not found", name)
}
