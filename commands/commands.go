package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config   string
	showVers bool

	version string
	commit  string

	// TractorCmd ...
	TractorCmd = &cobra.Command{
		Use:               "tractor",
		Short:             "tractor - data ingestion",
		Long:              ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: readConfig,
		PreRunE:           preFlight,
		RunE:              startTractor,
	}
)

func readConfig(ccmd *cobra.Command, args []string) error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if showVers {
		fmt.Printf("tractor %s (%s)\n", version, commit)
		return fmt.Errorf("")
	}

	return nil
}

func startTractor(ccmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	TractorCmd.AddCommand(printCmd)
	TractorCmd.AddCommand(runCmd)
	TractorCmd.Flags().BoolVarP(&showVers, "version", "v", false, "Display the current version of this CLI")

	viper.SetDefault("TRACTOR_PLUGINS_PATH", "bin/plugins")
	viper.SetDefault("TRACTOR_CHANNEL_BUFFER_SIZE", 1000)
	viper.SetDefault("TRACTOR_MAPPINGS_FILE", "mappings.yml")
}
