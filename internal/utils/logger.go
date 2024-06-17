package utils

import (
	"log"
	"os"
)

var (
	// logger provides a logging instance prefixed with "[http]" and standard flags.
	Logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

	// enable verbose logging.
	Verbose bool
)

// SetVerbose sets the verbosity level for logging
func SetVerbose(verbose bool) bool {
	Verbose = verbose
	if Verbose {
		Logger.SetFlags(log.LstdFlags | log.Lshortfile)
		Logger.Println("Verbose logging enabled")
	} else {
		Logger.SetFlags(log.LstdFlags)
	}
	return Verbose
}
