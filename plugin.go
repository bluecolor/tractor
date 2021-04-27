package tractor

import "github.com/bluecolor/tractor/config"

var Debug bool

type Initializer interface {
	Init(catalog *config.Catalog) error
}

type Validator interface {
	ValidateConfig() error
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

// Logger defines an plugin-related interface for logging.
type Logger interface {
	// Errorf logs an error message, patterned after log.Printf.
	Errorf(format string, args ...interface{})
	// Error logs an error message, patterned after log.Print.
	Error(args ...interface{})
	// Debugf logs a debug message, patterned after log.Printf.
	Debugf(format string, args ...interface{})
	// Debug logs a debug message, patterned after log.Print.
	Debug(args ...interface{})
	// Warnf logs a warning message, patterned after log.Printf.
	Warnf(format string, args ...interface{})
	// Warn logs a warning message, patterned after log.Print.
	Warn(args ...interface{})
	// Infof logs an information message, patterned after log.Printf.
	Infof(format string, args ...interface{})
	// Info logs an information message, patterned after log.Print.
	Info(args ...interface{})
}
