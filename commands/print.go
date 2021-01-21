package commands

import (
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/logging"
	"github.com/bluecolor/tractor/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	printInputPlugins  bool
	printOutputPlugins bool
)

func init() {
	printCmd.AddCommand(printMappingsCmd)
	printCmd.AddCommand(printPluginsCmd)

	printPluginsCmd.Flags().BoolVarP(&printInputPlugins, "input", "i", false, "Show input plugins")
	printPluginsCmd.Flags().BoolVarP(&printOutputPlugins, "output", "o", false, "Show output plugins")
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints the info requested with argumets",
	Run:   print,
}

var printMappingsCmd = &cobra.Command{
	Use:   "mappings",
	Short: "print mapping names to terminal",
	Run:   printMappings,
}

var printPluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "print available plugins in plugins path",
	Run:   printPlugins,
}

func print(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func printMappings(cmd *cobra.Command, args []string) {
	file := viper.GetString("TRACTOR_MAPPINGS_FILE")
	mappings := util.GetMappings(file)
	for i, mapping := range mappings {
		for name := range mapping {
			fmt.Printf("%6d %s\n", i+1, name)
		}
	}
	if len(mappings) == 0 {
		fmt.Printf("No mappings found\n")
	} else if len(mappings) == 1 {
		fmt.Printf("\nFound %d mapping\n", len(mappings))
	} else {
		fmt.Printf("\nFound %d mappings\n", len(mappings))
	}
}

func printPlugins(cmd *cobra.Command, args []string) {
	pluginsPath := viper.GetString("TRACTOR_PLUGINS_PATH")

	printInputPlugins = printInputPlugins || !(printInputPlugins || printOutputPlugins)
	printOutputPlugins = printOutputPlugins || !(printInputPlugins || printOutputPlugins)

	if printInputPlugins {
		plugins, err := util.GetPluginNamesByType(pluginsPath, "input")
		if err != nil {
			logging.Error("Failed to get input plugins")
		}
		fmt.Println("Input Plugins:")
		for i, plugin := range plugins {
			fmt.Printf("%4d %v\n", i+1, strings.TrimSuffix(plugin, ".so"))
		}
	}

	if printOutputPlugins {
		plugins, err := util.GetPluginNamesByType(pluginsPath, "output")
		if err != nil {
			logging.Error("Failed to get input plugins")
		}
		fmt.Println("Output Plugins:")
		for i, plugin := range plugins {
			fmt.Printf("%4d %v\n", i+1, strings.TrimSuffix(plugin, ".so"))
		}
	}
}
