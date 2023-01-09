package logger

import (
	"io"
	"log"
)

// NewDiscard returns a logger that discards all logs
func NewDiscard() *log.Logger {
	return log.New(io.Discard, "[all messages are discarded] ", 0)
}
