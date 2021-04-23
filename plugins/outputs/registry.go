package outputs

import "github.com/bluecolor/tractor"

type Creator func(config map[string]interface{}) tractor.Output

var Ouputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Ouputs[name] = creator
}
