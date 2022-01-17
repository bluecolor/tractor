package inputs

import (
	"github.com/bluecolor/tractor/lib/plugins"
	"github.com/bluecolor/tractor/lib/wire"
)

type InputPlugin interface {
	plugins.PluginDescriber
	Read(w wire.Wire) error
}
