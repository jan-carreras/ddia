package server

import (
	"ddia/src/resp"
	"strconv"
	"strings"
)

// Auth authenticates the client to the server, if requirepass directive is defined in the configuration file
// More: https://redis.io/commands/auth/
func (h *Handlers) Auth(c *client, expectedPassword string) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	clientPassword := c.args[1]

	if expectedPassword == "" {
		err := resp.NewError("ERR AUTH <password> called without any password configured for the default user. Are you sure your configuration is correct?")
		return c.writeResponse(err)
	}

	if clientPassword == expectedPassword {
		c.authenticated = true
		return c.writeResponse(resp.NewSimpleString("OK"))
	}

	err := resp.NewError("WRONGPASS invalid username-password pair or user is disabled.")
	return c.writeResponse(err)
}

// Echo returns message
//
//	redis> ECHO "Hello World!"
//	"Hello World!"
//
// More: https://redis.io/commands/echo/
func (h *Handlers) Echo(c *client) error {
	if len(c.args) <= 1 {
		return ErrWrongNumberArguments
	}

	return c.writeResponse(resp.NewStr(strings.Join(c.args[1:], " ")))
}

// Ping returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk
//
//	redis> PING
//	"PONG"
//	redis> PING "hello world"
//	"hello world"
//
// More: https://redis.io/commands/ping/
func (h *Handlers) Ping(c *client) error {
	if len(c.args) == 1 {
		return c.writeResponse(resp.NewSimpleString("PONG"))
	}

	response := resp.NewStr(strings.Join(c.args[1:], " "))
	return c.writeResponse(response)
}

// Quit asks the server to close the connection. The connection is closed as soon as all pending replies have been written to the client.
// More: https://redis.io/commands/quit/
func (h *Handlers) Quit(c *client) error {
	// TODO: Not quite it. We might we're writing a response somewhere else and we
	// should wait until we've finishing writing
	return c.conn.Close()
}

// Select the Redis logical database having the specified zero-based numeric index.
//
//	SELECT index
//
// More: https://redis.io/commands/select/
func (h *Handlers) Select(c *client, dbs []Storage) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	newDB := c.args[1]
	id, err := strconv.Atoi(newDB)
	if err != nil {
		return ErrValueNotInt
	} else if id < 0 || id >= len(dbs) {
		return ErrDBIndexOutOfRange
	}

	// Update the database where the client is pointing to
	c.dbIdx, c.db = id, dbs[id]

	return c.writeResponse(resp.NewSimpleString("OK"))
}
