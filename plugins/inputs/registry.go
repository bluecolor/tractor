package inputs

import (
	"github.com/bluecolor/tractor"
	cfg "github.com/bluecolor/tractor/config"
)

type Creator func(
	config map[string]interface{},
	catalog *cfg.Catalog,
	params map[string]interface{},
) (tractor.Input, error)

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
