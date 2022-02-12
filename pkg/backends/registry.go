package backends

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Creator func(config map[string]interface{}) (msg.FeedbackBackend, error)

var Backends = map[string]Creator{}

func Add(name string, creator Creator) {
	Backends[name] = creator
}

func GetBackend(name string, config map[string]interface{}) (msg.FeedbackBackend, error) {
	if creator, ok := Backends[name]; ok {
		return creator(config)
	}
	return nil, fmt.Errorf("connector %s not found", name)
}
