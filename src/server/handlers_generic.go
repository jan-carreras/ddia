package server

import (
	"ddia/src/resp"
	"errors"
	"strconv"
	"sync"
)

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
	err := h.atomic(c, func() error {
		for _, key := range keys {
			if c.db.Del(key) {
				countDeleted++
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(countDeleted))
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

	err := h.atomic(c, func() error {
		return c.db.Exists(key)
	})

	if errors.Is(err, ErrNotFound) {
		return c.writeResponse(resp.NewInteger(0))
	} else if err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(1))
}

// RandomKey return a random key from the currently selected database.
// More: https://redis.io/commands/randomkey/
func (h *Handlers) RandomKey(c *client) error {
	if err := c.requiredArgs(0); err != nil {
		return err
	}

	key, ok := "", false
	err := h.atomic(c, func() error {
		key, ok = c.db.RandomKey()
		return nil
	})

	if err != nil {
		return err
	}
	if !ok {
		return c.writeResponse(resp.NewNullStr())
	}

	return c.writeResponse(resp.NewStr(key))
}

// Rename renames key to newkey. It returns an error when key does not exist.
// More: https://redis.io/commands/rename/
func (h *Handlers) Rename(c *client) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, newKey := c.args[1], c.args[2]

	err := h.atomic(c, func() (err error) {
		return c.db.Rename(key, newKey)
	})

	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewStr("OK"))
}

// Move key from the currently selected database (see SELECT) to the specified
// destination database. When key already exists in the destination database, or
// it does not exist in the source database, it does nothing.
// More: https://redis.io/commands/move/
func (h *Handlers) Move(c *client, dbs []Storage, multiDBMutex *sync.Mutex) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, _dbIdx := c.args[1], c.args[2]

	dbIdx, err := strconv.Atoi(_dbIdx)
	if err != nil {
		return ErrValueNotInt
	}

	if dbIdx < 0 || dbIdx >= len(dbs) {
		return c.writeResponse(resp.NewInteger(0))
	}

	// This prevents death lock situation where:
	//   Process 1: locked DB 2, trying to acquire lock of DB 3
	//   Process 2: locked DB 3, trying to acquire lock of DB 2
	// If Process 1 and 2 need to fight for the same lock (multiDBMutex) then
	// this race condition disappear
	multiDBMutex.Lock()
	defer multiDBMutex.Unlock()

	err = h.atomic(c, func() error {
		v, err := c.db.Get(key)
		if err != nil {
			return err
		}

		otherDB := dbs[dbIdx]

		otherDB.Lock()
		defer otherDB.Unlock()

		if _, err := otherDB.Get(key); errors.Is(err, ErrNotFound) {
			// When key already exists in the destination database, [...] it does nothing.
		} else if err != nil {
			return err
		}

		if err := otherDB.Set(key, v); err != nil {
			return err
		}

		c.db.Del(key)

		return nil
	})

	if err != nil {
		return c.writeResponse(resp.NewInteger(0))
	}

	return c.writeResponse(resp.NewInteger(1))
}
