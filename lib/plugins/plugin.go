package plugins

import "github.com/bluecolor/tractor/lib/config"

type Initializer interface {
	Init() error
}

type Validator interface {
	Validate() error
}

type Discoverer interface {
	Discover() (*config.Catalog, error)
}

type Counter interface {
	Count() (int, error)
}

type PluginDescriber interface {
	SampleConfig() string

	Description() string
}
