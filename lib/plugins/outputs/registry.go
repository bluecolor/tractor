package outputs

import (
	"github.com/bluecolor/tractor/lib/config"
)

type Creator func(
	config map[string]interface{},
	sourceCatalog *config.Catalog,
	params map[string]interface{},
) (OutputPlugin, error)

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
