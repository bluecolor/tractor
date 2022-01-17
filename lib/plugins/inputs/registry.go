package inputs

import "github.com/bluecolor/tractor/lib/config"

type Creator func(
	config map[string]interface{},
	catalog *config.Catalog,
	params map[string]interface{},
) (InputPlugin, error)

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
