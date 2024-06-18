package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Logger provides a logging instance with JSON formatter.
	Logger = logrus.New()

	// Verbose indicates whether verbose logging is enabled.
	Verbose bool
)

func init() {
	// Set the output to stdout
	Logger.Out = os.Stdout

	// Set the default log level to Info
	Logger.SetLevel(logrus.InfoLevel)

	// Set the log format to JSON
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

// SetVerbose sets the verbosity level for logging
func SetVerbose(verbose bool) bool {
	Verbose = verbose
	if Verbose {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.Debug("Verbose logging enabled")
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
	return Verbose
}

// SetLogFormat allows overriding the logging format by setting custom log flags.
func SetLogFormat(formatter logrus.Formatter) {
	Logger.SetFormatter(formatter)
	Logger.Info("Log format overridden")
}
