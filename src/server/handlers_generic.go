package server

import (
	"ddia/src/expire"
	"ddia/src/resp"
	"errors"
	"strconv"
	"sync"
	"time"
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

func (h *Handlers) TTL(c *client, expire *expire.Expire) error {
	if err := c.requiredArgs(1); err != nil {
		return err
	}

	var ttl int

	key := c.args[1]
	if err := h.atomic(c, func() error {
		var ok bool

		if _, err := c.db.Get(key); errors.Is(err, ErrNotFound) {
			ttl = -2 // [...] if the key does not exist.
			return nil
		} else if err != nil {
			return err
		} else if ttl, ok = expire.TTL(key); !ok {
			ttl = -1 // [...] if the key exists but has no associated expire.
		}
		return nil
	}); err != nil {
		return nil
	}

	return c.writeResponse(resp.NewInteger(ttl))
}

// Expire set a timeout on key. After the timeout has expired, the key will
// automatically be deleted. A key with an associated timeout is often said to be
// volatile in Redis terminology.
//
// More: https://redis.io/commands/expire/
func (h *Handlers) Expire(c *client, expire *expire.Expire) error {
	if err := c.requiredArgs(2); err != nil {
		return err
	}

	key, _seconds := c.args[1], c.args[2]

	seconds, err := strconv.Atoi(_seconds)
	if err != nil {
		return ErrWrongKind
	}

	result := 1 // if the timeout was set.

	expirationTime := time.Now().Add(time.Duration(seconds) * time.Second).Unix()
	if err := h.atomic(c, func() error {
		// Make sure that
		_, err := c.db.Get(key)
		if errors.Is(err, ErrNotFound) {
			result = 0 // if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments.
			return nil // No need
		} else if err != nil {
			return err
		}

		expire.AddUpdate(c.dbIdx, key, expirationTime)
		return nil
	}); err != nil {
		return err
	}

	return c.writeResponse(resp.NewInteger(result))
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
