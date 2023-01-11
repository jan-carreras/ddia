package server

import (
	"fmt"
	"io"
	"net"
)

type client struct {
	// conn is the TCP socket that connects with the client
	conn net.Conn
	// args are the commands being sent by the network
	args []string
	// dbIdx is the ID of the database where the client is connected to. Default to DB 0
	dbIdx int
	// db points to the active database the client is connected to. Can be changed used the SELECT command
	db Storage
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

// command returns the command name (c.args[0]), or empty
func (c *client) command() string {
	if len(c.args) == 0 {
		return ""
	}
	return c.args[0]
}

// writeResponse writes into the active connection, returning an error if it fails
func (c *client) writeResponse(to io.WriterTo) error {
	if _, err := to.WriteTo(c.conn); err != nil {
		return fmt.Errorf("unable to writeResponse to the client: %w", err)
	}

	return nil
}
