package log

import (
	"io"
	"log"
)

func ServerLogger() *log.Logger {
	return log.New(loggerOutput(), "[server] ", 0)
}

func ClientLogger() *log.Logger {
	return log.New(loggerOutput(), "[client] ", 0)
}

func loggerOutput() io.Writer {
	return io.Discard
}
