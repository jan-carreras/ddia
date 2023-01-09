// Package log is used to get loggers for testing purposes
package log

import (
	"io"
	"log"
)

// ServerLogger returns a logger to be used with the Server
func ServerLogger() *log.Logger {
	return log.New(loggerOutput(), "[server] ", 0)
}

// ClientLogger returns a logger to be used with the Client
func ClientLogger() *log.Logger {
	return log.New(loggerOutput(), "[client] ", 0)
}

func loggerOutput() io.Writer {
	return io.Discard
}
