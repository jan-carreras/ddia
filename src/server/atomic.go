package server

import (
	"bufio"
	"ddia/src/resp"
	"strconv"
)

// atomic runs fnx on the selected databases by the client acquiring the DB lock
// first and releasing it safely.
// If the change has been done correctly in memory, we write the operation on the AOF file
// with the lock still acquired.
// Atomic does not support atomic operations between two different databases.
func (h *Handlers) atomic(c *client, fnx func() error) error {
	c.db.Lock()
	defer c.db.Unlock()

	if err := fnx(); err != nil {
		return err
	}

	if err := h.writeToAOF(c); err != nil {
		return err
	}

	return nil
}

// writeToAOF persists the executed command if the AOF storage has been set, and
// the command being executed is a "write" command. It always pre-appends the
// SELECT {DB_ID} number before each command, to make sure that operation is
// going to be re-played in the correct DB. It obvious that we could memorize
// into which DB did we write the last time, and avoid the same SELECT over and
// over. It's an optimization to be done in the future.
func (h *Handlers) writeToAOF(c *client) error {
	if h.aof == nil {
		return nil
	}

	cmd, ok := command(c.command())

	if !ok {
		return nil
	}

	if cmd.Operation != "write" {
		return nil
	}

	buf := bufio.NewWriter(h.aof)

	sel := resp.NewArray([]string{"SELECT", strconv.Itoa(c.dbIdx)})
	if _, err := sel.WriteTo(buf); err != nil {
		return err
	}

	if _, err := c.argsWriter.WriteTo(buf); err != nil {
		return err
	}

	if err := buf.Flush(); err != nil {
		return err
	}

	return nil
}
