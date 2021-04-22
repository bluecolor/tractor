package main

import (
	"fmt"

	"github.com/bluecolor/tractor/cmd"
	_ "github.com/bluecolor/tractor/plugins/inputs/all"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func setupLogger() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := cmd.TractorCmd.PersistentFlags().Lookup("loglevel").Value.String()
	if logLevel == "" {
		logLevel = viper.GetString("TRACTOR_LOG_LEVEL")
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	return nil
}

func main() {
	godotenv.Load()
	setupLogger()
	err := cmd.TractorCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
