package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"fmt"
	"io"
)

// Handlers define the commands being handled by the Redis Server. A new command should be registered
// as a public function in the handler file.
type Handlers struct {
	logger logger.Logger
	aof    io.Writer
}

// NewHandlers returns a Handlers
func NewHandlers(logger logger.Logger, aof io.Writer) *Handlers {
	return &Handlers{logger: logger, aof: aof}
}

// UnknownCommand returns an error when the command is unknown
func (h *Handlers) UnknownCommand(c *client) error {
	err := resp.NewError(fmt.Sprintf("ERR unknown command '%s'", c.command()))
	return c.writeResponse(err)
}
