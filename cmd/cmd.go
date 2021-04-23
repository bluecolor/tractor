package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config      string
	showVersion bool
	logLevel    string
	mapping     string
	inputFlag   bool
	outputFlag  bool

	version string
	commit  string

	// TractorCmd ...
	TractorCmd = &cobra.Command{
		Use:               "tractor",
		Short:             "🚜 tractor - data ingestion tool",
		Long:              ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: readConfig,
		PreRunE:           preFlight,
		RunE:              start,
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
	if showVersion {
		fmt.Printf("tractor %s (%s)\n", version, commit)
		return fmt.Errorf("")
	}

	return nil
}

func start(ccmd *cobra.Command, args []string) error {

	return nil
}

func init() {
	viper.SetDefault("TRACTOR_CHANNEL_BUFFER_SIZE", 1000)
	viper.SetDefault("TRACTOR_MAPPINGS_FILE", "mappings.yml")
	viper.SetDefault("TRACTOR_LOG_LEVEL", "info")

	TractorCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display the current version of this CLI")
	TractorCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "", "log level")
	TractorCmd.PersistentFlags().StringVar(&config, "config", "tractor.yml", "Config file")

	TractorCmd.AddCommand(runCmd)
	TractorCmd.AddCommand(pluginCmd)
}
