package cmd

import (
	"fmt"
	"os"

	"github.com/bluecolor/tractor"
	cfg "github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config      string
	showVersion bool
	logLevel    string
	mapping     string

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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run given mapping",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	if mapping == "" {
		println("Mapping is not given")
		os.Exit(1)
	}
	conf := cfg.NewConfig()
	err := conf.LoadConfig(config)
	if err != nil {
		println("Failed to load config file")
		os.Exit(1)
	}
	m, err := conf.GetMapping(mapping)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if creator, ok := inputs.Inputs[m.Input.Plugin]; ok {
		inputPlugin := creator(m.Input.Config)
		if validator, ok := inputPlugin.(tractor.Validator); ok {
			if err := validator.ValidateConfig(); err != nil {
				println("❌  Failed to validate input config")
				println(err.Error())
				os.Exit(1)
			} else {
				println("✔️  Input config validated")
			}
		}
	} else {
		println(fmt.Sprintf("❌  No matching plugin found for %s", m.Input.Plugin))
		os.Exit(1)
	}
}

func start(ccmd *cobra.Command, args []string) error {
	// conf := cfg.NewConfig()
	// err := conf.LoadConfig(config)
	// if err != nil {
	// 	return err
	// }
	// if err != nil {
	// 	return err
	// }
	// for k := range inputs.Inputs {
	// 	println(k)
	// }
	return nil
}

func init() {
	viper.SetDefault("TRACTOR_CHANNEL_BUFFER_SIZE", 1000)
	viper.SetDefault("TRACTOR_MAPPINGS_FILE", "mappings.yml")
	viper.SetDefault("TRACTOR_LOG_LEVEL", "info")

	TractorCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Display the current version of this CLI")
	TractorCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "", "log level")
	TractorCmd.PersistentFlags().StringVar(&config, "config", "tractor.yml", "config file")
	TractorCmd.PersistentFlags().StringVar(&mapping, "mapping", "", "mapping name")
	TractorCmd.AddCommand(runCmd)
}
