package server

import (
	"ddia/src/resp"
	"ddia/src/server/config"
	"fmt"
	"strings"
)

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
		var value *resp.Array
		switch key {
		case "save":
			value = resp.NewArray([]string{"save", "3600 1 300 100 60 10000"})
		case "appendonly":
			value = resp.NewArray([]string{"appendonly", "no"})
		default:
			v, _ := config.Get(key)
			value = resp.NewArray([]string{key, v})
		}

		return c.writeResponse(value)
	}

	err := resp.NewError(fmt.Sprintf("ERR unknown subcommand '%s'.", cmd))
	return c.writeResponse(err)
}

// DBSize : Return the number of keys in the selected database
// More: https://redis.io/commands/dbsize/
func (h *Handlers) DBSize(c *client) error {
	if err := c.requiredArgs(0); err != nil {
		return nil
	}

	return c.writeResponse(resp.NewInteger(c.db.Size()))
}

// FlushAll delete all the keys of all the existing databases, not just the currently selected one.
// More: https://redis.io/commands/flushall
func (h *Handlers) FlushAll(c *client, dbs []Storage) error {
	if err := c.requiredArgs(0); err != nil {
		return err
	}

	for _, db := range dbs {
		db.Lock()
		if err := db.FlushDB(); err != nil {
			db.Unlock()
			return err
		}
		db.Unlock()
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
}

// FlushDB delete all the keys of the currently selected DB. This command never fails.
// More: https://redis.io/commands/flushdb/
func (h *Handlers) FlushDB(c *client) error {
	if err := c.requiredArgs(0); err != nil {
		return err
	}

	err := h.atomic(c, func() error {
		return c.db.FlushDB()
	})
	if err != nil {
		return err
	}

	return c.writeResponse(resp.NewSimpleString("OK"))
}
