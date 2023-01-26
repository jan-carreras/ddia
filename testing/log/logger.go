// Package log is used to get loggers for testing purposes
package log

import (
	"io"
	"log"
)

// ServerLogger returns a logger to be used with the Server
func ServerLogger() *log.Logger {
	return log.New(io.Discard, "[server] ", 0)
}
