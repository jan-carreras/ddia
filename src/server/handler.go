package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"fmt"
	"strconv"
	"strings"
)

// Handlers define the commands being handled by the Redis Server. A new command should be registered
// as a public function in the handler file.
type Handlers struct {
	logger  logger.Logger
	storage Storage
}

// NewHandlers returns a Handlers
func NewHandlers(logger logger.Logger, storage Storage) *Handlers {
	return &Handlers{logger: logger, storage: storage}
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

	value, err := h.storage.Get(key)
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

	err := h.storage.Set(key, value)
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

	newValue, err := h.storage.IncrementBy(key, incr)
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

	newValue, err := h.storage.Increment(key)
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

	newValue, err := h.storage.DecrementBy(key, decrement)
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

	newV, err := h.storage.Decrement(key)
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

	return c.writeResponse(resp.NewInteger(strconv.Itoa(h.storage.Size())))
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
		if h.storage.Del(key) {
			countDeleted++
		}
	}

	return c.writeResponse(resp.NewInteger(strconv.Itoa(countDeleted)))
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

// UnknownCommand returns an error when the command is unknown
func (h *Handlers) UnknownCommand(c *client) error {
	err := resp.NewError(fmt.Sprintf("ERR unknown command '%s'", c.command()))
	return c.writeResponse(err)
}
