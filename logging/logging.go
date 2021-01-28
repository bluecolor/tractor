package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//AppEnv for config
	AppEnv = "APP_ENV"
	log    *zap.SugaredLogger
)

func init() {
	BuildLogger(os.Getenv(AppEnv))
}

// BuildLogger builds log config
func BuildLogger(env string) {
	var outputPaths []string
	var level zapcore.Level

	if env == "development" || env == "" {
		outputPaths = []string{"stdout"}
		level = zapcore.DebugLevel
	} else if env == "production" {
		outputPaths = []string{"./tractor.log"}
		level = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Info ...
func Info(args ...interface{}) {
	log.Info(args...)
}
