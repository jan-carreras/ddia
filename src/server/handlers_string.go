package server

import (
	"ddia/src/resp"
	"errors"
	"strconv"
)

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

	return h.incrBy(c, key, -1)
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

	return h.incrBy(c, key, -decrement)
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

	var value string
	err := h.atomic(c, func() (err error) {
		value, err = c.db.Get(key)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr(value))
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

	return h.incrBy(c, key, incr)
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

	return h.incrBy(c, key, 1)
}

func (h *Handlers) incrBy(c *client, key string, incr int) error {
	var newValue string
	err := h.atomic(c, func() (err error) {
		newValue, err = c.db.IncrementBy(key, incr)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString(newValue))
}

// MGet returns the values of all specified keys.
// redis> SET key1 "Hello"
// "OK"
// redis> SET key2 "World"
// "OK"
// redis> MGET key1 key2 nonexisting
// 1) "Hello"
// 2) "World"
// 3) (nil)
// More: https://redis.io/commands/mget/
func (h *Handlers) MGet(c *client) error {
	if len(c.args) <= 1 {
		return ErrWrongNumberArguments
	}

	keys := c.args[1:]

	values := make([]string, 0)
	err := h.atomic(c, func() error {
		for _, key := range keys {
			value, err := c.db.Get(key)
			if err != nil {
				return err
			}
			values = append(values, value)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewArray(values))
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

	err := h.atomic(c, func() error {
		return c.db.Set(key, value)
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
}

// SetNX set key to hold string value if key does not exist.
// redis> SETNX mykey "Hello"
// (integer) 1
// redis> SETNX mykey "World"
// (integer) 0
// redis> GET mykey
// "Hello"
// More: https://redis.io/commands/setnx/
func (h *Handlers) SetNX(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, value := c.args[1], c.args[2]

	response := 0
	err := h.atomic(c, func() error {
		err := c.db.Exists(key)
		if errors.Is(err, ErrNotFound) {
			response = 1
			return c.db.Set(key, value)
		}
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(response))
}

// Substr returns the substring of the string value stored at key, determined by
// the offsets start and end (both are inclusive)
// redis> SET mykey "This is a string"
// "OK"
// redis> GETRANGE mykey 0 3
// "This"
// redis> GETRANGE mykey -3 -1
// "ing"
// redis> GETRANGE mykey 0 -1
// "This is a string"
// redis> GETRANGE mykey 10 100
// "string"
// Note: https://redis.io/commands/substr/
func (h *Handlers) Substr(c *client) error {
	if err := c.requiredArgs(3); err != nil {
		return err
	}

	key := c.args[1]

	start, err := strconv.Atoi(c.args[2])
	if err != nil {
		return ErrValueNotInt
	}

	end, err := strconv.Atoi(c.args[3])
	if err != nil {
		return ErrValueNotInt
	}

	var value string
	err = h.atomic(c, func() (err error) {
		value, err = c.db.Get(key)
		return err
	})

	if err != nil {
		return err
	}

	// start is out of bounds
	if start > end || start > len(value) {
		return c.writeResponse(resp.NewStr(""))
	}

	// the offsets start and end (both are inclusive), thus the +1 on the end
	end++

	// Prevent out of bounds offset
	if end >= len(value) {
		end = len(value)
	}

	return c.writeResponse(resp.NewStr(value[start:end]))
}
