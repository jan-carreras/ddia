package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"ddia/src/server/config"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Handlers define the commands being handled by the Redis Server. A new command should be registered
// as a public function in the handler file.
type Handlers struct {
	logger logger.Logger
}

// NewHandlers returns a Handlers
func NewHandlers(logger logger.Logger) *Handlers {
	return &Handlers{logger: logger}
}

// Get the value of a key
//
//	GET key
//
// More: https://redis.io/commands/get/
func (h *Handlers) Get(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	value, err := c.db.Get(key)
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString(value))
}

// Set key to hold the string value
//
//	SET key value
//
// More: https://redis.io/commands/set/
func (h *Handlers) Set(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, value := c.args[1], c.args[2]

	err := c.db.Set(key, value)
	if err != nil {
		return fmt.Errorf("storage.Set: %w", err)
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
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

	response := resp.NewSimpleString(strings.Join(c.args[1:], " "))
	return c.writeResponse(response)

}

// IncrBy increments the number stored at key by increment.
//
//	INCRBY key increment
//
// More: https://redis.io/commands/incrby/
func (h *Handlers) IncrBy(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, increment := c.args[1], c.args[2]

	incr, err := strconv.Atoi(increment)
	if err != nil {
		return ErrValueNotInt
	}

	newValue, err := c.db.IncrementBy(key, incr)
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString(newValue))
}

// Incr increments the number stored at key by one.
//
//	INCR key
//
// More: https://redis.io/commands/incr/
func (h *Handlers) Incr(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	newValue, err := c.db.Increment(key)
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString(newValue))
}

// DecrBy decrements the number stored at key by decrement.
//
//	DECRBY key decrement
//
// More: https://redis.io/commands/decrby/
func (h *Handlers) DecrBy(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, value := c.args[1], c.args[2]

	decrement, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	newValue, err := c.db.DecrementBy(key, decrement)
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString(newValue))
}

// Decr decrements the number stored at key by one.
//
//	DECR key
//
// More: https://redis.io/commands/decr/
func (h *Handlers) Decr(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	newV, err := c.db.Decrement(key)
	if err != nil {
		return nil
	}

	return c.writeResponse(resp.NewSimpleString(newV))
}

// DBSize : Return the number of keys in the selected database
// More: https://redis.io/commands/dbsize/
func (h *Handlers) DBSize(c *client) error {
	if err := c.requiredArgs(0); err != nil {
		return nil
	}

	return c.writeResponse(resp.NewInteger(c.db.Size()))
}

// Del removes the specified keys. A key is ignored if it does not exist.
//
//	redis> SET key1 "Hello"
//	"OK"
//	redis> SET key2 "World"
//	"OK"
//	redis> DEL key1 key2 key3
//	(integer) 2
//
// More: https://redis.io/commands/del/
func (h *Handlers) Del(c *client) error {
	if len(c.args) <= 1 {
		return ErrWrongNumberArguments
	}

	keys := c.args[1:]

	countDeleted := 0
	for _, key := range keys {
		if c.db.Del(key) {
			countDeleted++
		}
	}

	return c.writeResponse(resp.NewInteger(countDeleted))
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

	return c.writeResponse(resp.NewSimpleString(strings.Join(c.args[1:], " ")))
}

// Quit asks the server to close the connection. The connection is closed as soon as all pending replies have been written to the client.
// More: https://redis.io/commands/quit/
func (h *Handlers) Quit(c *client) error {
	// TODO: Not quite it. We might we writing a response somewhere else and we
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

// UnknownCommand returns an error when the command is unknown
func (h *Handlers) UnknownCommand(c *client) error {
	err := resp.NewError(fmt.Sprintf("ERR unknown command '%s'", c.command()))
	return c.writeResponse(err)
}

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

// FlushDB delete all the keys of the currently selected DB. This command never fails.
// More: https://redis.io/commands/flushdb/
func (h *Handlers) FlushDB(c *client) error {
	if err := c.requiredArgs(0); err != nil {
		return err
	}

	if err := c.db.FlushDB(); err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
}

// FlushAll delete all the keys of all the existing databases, not just the currently selected one.
// More: https://redis.io/commands/flushall
func (h *Handlers) FlushAll(c *client, dbs []Storage) error {
	if err := c.requiredArgs(0); err != nil {
		return err
	}

	for _, db := range dbs {
		if err := db.FlushDB(); err != nil {
			return err
		}
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
}

// Exists returns if key exists. 1 if exists, 0 otherwiese.
//
//	redis> SET key1 "Hello"
//	"OK"
//	redis> EXISTS key1
//	(integer) 1
//	redis> EXISTS nosuchkey
//	(integer) 0
//
// More: https://redis.io/commands/exists/
func (h *Handlers) Exists(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	err := c.db.Exists(key)
	if errors.Is(err, ErrNotFound) {
		return c.writeResponse(resp.NewInteger(0))
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(1))
}

// Config returns stuff from the Config.
//
// TODO: This command has been included to try to make redis-benchmark cli to work. I'm returning hardcoded stuff
// in the hope that the command will work. Without this there is no hope
func (h *Handlers) Config(c *client, config config.Config) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	cmd := c.args[1]
	if strings.ToUpper(cmd) == "GET" {
		key := c.args[2]
		value := resp.Array{}
		switch key {
		case "save":
			value = resp.NewArray([]string{"save", "3600 1 300 100 60 10000"})
		case "appendonly":
			value = resp.NewArray([]string{"appendonly", "no"})
		default:
			v, _ := config.Get(key)
			value = resp.NewArray([]string{key, v})
		}

		return c.writeResponse(&value)
	}

	err := resp.NewError(fmt.Sprintf("ERR unknown subcommand '%s'.", cmd))
	return c.writeResponse(err)
}
