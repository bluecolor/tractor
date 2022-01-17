package outputs

import (
	"github.com/bluecolor/tractor/lib/plugins"
	"github.com/bluecolor/tractor/lib/wire"
)

type OutputPlugin interface {
	plugins.PluginDescriber
	Write(w *wire.Wire) error
}
