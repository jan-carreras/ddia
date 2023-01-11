package server

import (
	"fmt"
	"io"
	"net"
)

type client struct {
	conn net.Conn
	args []string
}

// requiredArgs makes sure that the number of arguments is equal to expectedArguments. Returns error otherwise.
// Note that c.args[0] is the command, and this is not taken into account in the expectedArguments
// Example: client{args: []string{"set", "hello", "world"}.requiredArgs(2) == true
func (c *client) requiredArgs(expectedArguments int) error {
	expectedArguments++ // c.args[0] is always the command
	if len(c.args) != expectedArguments {
		return ErrWrongNumberArguments
	}

	return nil
}

func (c *client) command() string {
	if len(c.args) == 0 {
		return ""
	}
	return c.args[0]
}

func (c *client) writeResponse(to io.WriterTo) error {
	if _, err := to.WriteTo(c.conn); err != nil {
		return fmt.Errorf("unable to writeResponse to the client: %w", err)
	}

	return nil
}
