package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"errors"
	"fmt"
	"net"
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

// Get : Get the value of a key
func (h *Handlers) Get(conn net.Conn, cmd []string) error {
	if len(cmd) != 2 {
		err := resp.NewError("ERR wrong number of arguments for 'GET' command")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	val, err := h.storage.Get(cmd[1])
	if errors.Is(err, ErrNotFound) {
		ok := resp.NewSimpleString("")
		if _, err := ok.WriteTo(conn); err != nil {
			h.logger.Printf("unable to write")
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("storage.Get: %w", err)
	}

	ok := resp.NewSimpleString(val)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write")
	}
	return nil
}

// Set : Set the string value of a key
func (h *Handlers) Set(conn net.Conn, cmd []string) error {
	if len(cmd) != 3 {
		err := resp.NewError("ERR wrong number of arguments for 'SET' command")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}
	err := h.storage.Set(cmd[1], cmd[2])
	if err != nil {
		return fmt.Errorf("storage.Set: %w", err)
	}

	ok := resp.NewSimpleString("OK")
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write")
	}

	return nil
}

// UnknownCommand returns an error when the command is unknown
func (h *Handlers) UnknownCommand(conn net.Conn, verb string) error {
	err := resp.NewError(fmt.Sprintf("ERR unknown command '%s'", verb))

	if _, err := err.WriteTo(conn); err != nil {
		return fmt.Errorf("unable to write command not found: %w", err)
	}

	return nil
}

// Ping : Ping the server
func (h *Handlers) Ping(conn net.Conn, cmd []string) error {
	if len(cmd) == 1 {
		ok := resp.NewSimpleString("PONG")
		if _, err := ok.WriteTo(conn); err != nil {
			h.logger.Printf("unable to write: %v", err)
		}

		return nil
	}

	ok := resp.NewSimpleString(strings.Join(cmd[1:], " "))
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// IncrBy increments the number stored at key by increment.
//
//	INCRBY key increment
//
// More: https://redis.io/commands/incrby/
func (h *Handlers) IncrBy(conn net.Conn, cmd []string) error {
	if len(cmd) != 3 {
		err := resp.NewErrorWrongArguments(cmd[0])
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	key, value := cmd[1], cmd[2]

	v, err := strconv.Atoi(value)
	if err != nil {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	newV, err := h.storage.IncrementBy(key, v)
	if errors.Is(err, ErrValueNotInt) || errors.Is(err, ErrWrongKind) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	ok := resp.NewSimpleString(newV)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// Incr increments the number stored at key by one.
//
//	INCR key
//
// More: https://redis.io/commands/incr/
func (h *Handlers) Incr(conn net.Conn, cmd []string) error {
	if len(cmd) != 2 {
		err := resp.NewErrorWrongArguments(cmd[0])
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	key := cmd[1]

	newV, err := h.storage.Increment(key)
	if errors.Is(err, ErrValueNotInt) || errors.Is(err, ErrWrongKind) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	ok := resp.NewSimpleString(newV)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// DecrBy decrements the number stored at key by decrement.
//
//	DECRBY key decrement
//
// More: https://redis.io/commands/decrby/
func (h *Handlers) DecrBy(conn net.Conn, cmd []string) error {
	if len(cmd) != 3 {
		err := resp.NewErrorWrongArguments(cmd[0])
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	key, value := cmd[1], cmd[2]

	v, err := strconv.Atoi(value)
	if err != nil {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	newV, err := h.storage.DecrementBy(key, v)
	if errors.Is(err, ErrValueNotInt) || errors.Is(err, ErrWrongKind) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	ok := resp.NewSimpleString(newV)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// Decr decrements the number stored at key by one.
//
//	DECR key
//
// More: https://redis.io/commands/decr/
func (h *Handlers) Decr(conn net.Conn, cmd []string) error {
	if len(cmd) != 2 {
		err := resp.NewErrorWrongArguments(cmd[0])
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	key := cmd[1]

	newV, err := h.storage.Decrement(key)
	if errors.Is(err, ErrValueNotInt) || errors.Is(err, ErrWrongKind) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	ok := resp.NewSimpleString(newV)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// DBSize : Return the number of keys in the selected database
func (h *Handlers) DBSize(conn net.Conn, cmd []string) error {
	size := resp.NewInteger(strconv.Itoa(h.storage.Size()))
	if _, err := size.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write: %v", err)
	}

	return nil
}

// Del : Delete a key
func (h *Handlers) Del(conn net.Conn, cmd []string) error {
	if len(cmd) == 1 {
		err := resp.NewError("ERR wrong number of arguments for 'del' command")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("unable to write all the configuration: %w", err)
		}
		return nil
	}

	countDeleted := 0
	for _, key := range cmd[1:] {
		if h.storage.Del(key) {
			countDeleted++
		}
	}

	r := resp.NewInteger(strconv.Itoa(countDeleted))
	if _, err := r.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write response: %v", err)
	}

	return nil
}
