package cmd

import (
	"fmt"
	"os"

	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Plugin info",
	Run:   plugin,
}

func plugin(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("❌  Missing argument for 'plugin' command")
		os.Exit(1)
	}
	i, o := inputFlag, outputFlag
	if !i && !o {
		i, o = true, true
	}
	if i {
		if creator, ok := inputs.Inputs[args[0]]; ok {
			inputPlugin, err := creator(make(map[string]interface{}), nil, nil)
			if err != nil {
				println("Failed to create plugin")
				os.Exit(1)
			}
			println(fmt.Sprintf("🧩 %s - Input plugin sample config", args[0]))
			println(inputPlugin.SampleConfig())
		} else {
			println(fmt.Sprintf("⚠️  No matching input plugin found for %s", args[0]))
		}
	}
	if o {
		// todo
	}

}

func init() {
	pluginCmd.PersistentFlags().BoolVar(&inputFlag, "input", false, "Input flag")
	pluginCmd.PersistentFlags().BoolVar(&outputFlag, "output", false, "Output flag")
}
