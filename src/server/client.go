package server

import (
	"bufio"
	"ddia/src/resp"
	"fmt"
	"io"
)

type client struct {
	// conn is the TCP socket that connects with the client
	conn io.ReadWriteCloser
	// reader is a wrapper for conn using bufio
	reader *bufio.Reader

	// args are the commands being sent by the network
	args []string
	// argsWriter allows to write the command we've received into a writer
	argsWriter io.WriterTo

	// dbIdx is the ID of the database where the client is connected to. Default to DB 0
	dbIdx int
	// db points to the active database the client is connected to. Can be changed used the SELECT command
	db Storage
	// authenticated is true when the client has successfully authenticated to the Server using the AUTH command
	authenticated bool
}

// newClient returns a client
func newClient(conn io.ReadWriteCloser, db Storage) *client {
	return &client{conn: conn, db: db, reader: bufio.NewReader(conn)}
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

func (c *client) readCommand() error {
	operation, err := resp.ReadOperation(c.reader)
	if err != nil {
		return fmt.Errorf("unable to peak operation: %w", err)
	}

	switch operation {
	case resp.ArrayOp:
		args, err := c.parseBulkString(c.reader)
		if err != nil {
			return fmt.Errorf("parsing Bulk String command: %w", err)
		}

		// Load the arguments to the client, to be able to process the request
		c.args = args.Strings()
		c.argsWriter = args

		return nil
	case resp.RawPing: // Ping, but without being part of SimpleString. I don't know which part of the specs describes this :/
		var s resp.SimpleString
		_, err := s.ReadFrom(c.reader)
		cmd := "P" + s.String()
		if cmd == "PING" {
			c.args = []string{cmd}
		}
		return err
	default:
		return fmt.Errorf("unknown opertion type: %q", operation)
	}
}

func (c *client) parseBulkString(conn io.Reader) (*resp.Array, error) {
	b := &resp.Array{}
	_, err := b.ReadFrom(conn)
	return b, err
}

func (c *client) close() error {
	return c.conn.Close()
}
