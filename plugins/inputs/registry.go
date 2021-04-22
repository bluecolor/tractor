package inputs

import "github.com/bluecolor/tractor"

type Creator func(config map[string]interface{}) tractor.Input

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
