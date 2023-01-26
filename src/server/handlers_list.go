package server

import (
	"ddia/src/resp"
	"errors"
	"strconv"
)

// LRange Returns the specified elements of the list stored at key. The offsets
// start and stop are zero-based indexes, with 0 being the first element of the
// list (the head of the list), 1 being the next element and so on. These offsets
// can also be negative numbers indicating offsets starting at the end of the
// list. For example, -1 is the last element of the list, -2 the penultimate, and
// so on.
// More: https://redis.io/commands/lrange/
//
// TODO: Optimisation, we don't necessarily need to start from the Start of the list. We can traverse it
// starting from the end, if the requested range is "nearer" the end.
func (h *Handlers) LRange(c *client) error {
	if err := c.requiredArgs(3); err != nil {
		return err
	}

	key, _start, _stop := c.args[1], c.args[2], c.args[3]

	start, err := strconv.Atoi(_start)
	if err != nil {
		return err
	}

	stop, err := strconv.Atoi(_stop)
	if err != nil {
		return err
	}

	var values []string
	err = h.atomic(c, func() (err error) {
		values, err = c.db.LRange(key, start, stop)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewArray(values))
}

// LRem removes the first count occurrences of elements equal to element from the
// list stored at key. More: https://redis.io/commands/lrem/
func (h *Handlers) LRem(c *client) error {
	if err := c.requiredArgs(3); err != nil {
		return err
	}

	key, _count, element := c.args[1], c.args[2], c.args[3]

	count, err := strconv.Atoi(_count)
	if err != nil {
		return err
	}

	var deleted int
	err = h.atomic(c, func() error {
		deleted, err = c.db.LRem(key, count, element)
		return err
	})

	if errors.Is(err, ErrNotFound) {
		// If the list does not exist, just report that 0 records have been removed
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(deleted))
}

// LIndex returns the element at "index" index in the list stored at key. The index
// is zero-based, so 0 means the first element, 1 the second element and so on.
// Negative indices can be used to designate elements starting at the tail of the
// list. Here, -1 means the last element, -2 means the penultimate and so forth.
// More: https://redis.io/commands/lindex/
func (h *Handlers) LIndex(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, _index := c.args[1], c.args[2]

	index, err := strconv.Atoi(_index)
	if err != nil {
		return err
	}

	var value string
	err = h.atomic(c, func() (err error) {
		value, err = c.db.LIndex(key, index)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr(value))
}

// LSet sets the list element at index to element.
// More: https://redis.io/commands/lset/
func (h *Handlers) LSet(c *client) error {
	if err := c.requiredArgs(3); err != nil {
		return err
	}

	key, _index, element := c.args[1], c.args[2], c.args[3]

	index, err := strconv.Atoi(_index)
	if err != nil {
		return ErrValueNotInt
	}

	err = h.atomic(c, func() error {
		return c.db.LSet(key, index, element)
	})

	if errors.Is(err, ErrNotFound) {
		return c.writeResponse(resp.NewNullSimpleString())
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr("OK"))
}

// LLen returns the length of the list stored at key
// More: https://redis.io/commands/llen/
func (h *Handlers) LLen(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	var length int
	err := h.atomic(c, func() (err error) {
		length, err = c.db.LLen(key)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(length))
}

// LPush insert all the specified values at the head of the list stored at key.
// If key does not exist, it is created as empty list before performing the push
// operations. When key holds a value that is not a list, an error is returned.
// More: https://redis.io/commands/lpush/
func (h *Handlers) LPush(c *client) error {
	if len(c.args) < 2 {
		return ErrWrongNumberArguments
	}

	key, elements := c.args[1], c.args[2:]

	var n int
	err := h.atomic(c, func() (err error) {
		n, err = c.db.LPush(key, elements)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(n))
}

// RPush insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push
// operation. When key holds a value that is not a list, an error is returned.
// More: https://redis.io/commands/rpush/
func (h *Handlers) RPush(c *client) error {
	if len(c.args) < 2 {
		return ErrWrongNumberArguments
	}

	key, elements := c.args[1], c.args[2:]

	var n int
	err := h.atomic(c, func() (err error) {
		n, err = c.db.RPush(key, elements)
		return err
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(n))
}

// LPop removes and returns the first elements of the list stored at key.
// More: https://redis.io/commands/lpop/
func (h *Handlers) LPop(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	var value string
	err := h.atomic(c, func() (err error) {
		value, err = c.db.LPop(key)
		return err
	})

	if errors.Is(err, ErrNotFound) {
		return c.writeResponse(resp.NewNullStr())
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr(value))
}

// RPop removes and returns the last elements of the list stored at key.
// More: https://redis.io/commands/rpop/
func (h *Handlers) RPop(c *client) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	key := c.args[1]

	var value string
	err := h.atomic(c, func() (err error) {
		value, err = c.db.RPop(key)
		return err
	})

	if errors.Is(err, ErrNotFound) {
		return c.writeResponse(resp.NewNullStr())
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr(value))
}
