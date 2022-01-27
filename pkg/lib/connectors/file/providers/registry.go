package providers

import "fmt"

type Creator func(config map[string]interface{}) (Provider, error)

var Providers = map[string]Creator{}

func Add(name string, creator Creator) {
	Providers[name] = creator
}

func GetProvider(name string, config map[string]interface{}) (Provider, error) {
	if creator, ok := Providers[name]; ok {
		return creator(config)
	}
	return nil, fmt.Errorf("provider %s not found", name)
}
