package outputs

import (
	"github.com/bluecolor/tractor"
	cfg "github.com/bluecolor/tractor/config"
)

type Creator func(
	config map[string]interface{},
	catalog *cfg.Catalog,
	params map[string]interface{},
) (tractor.Output, error)

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
