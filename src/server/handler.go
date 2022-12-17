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

type Handlers struct {
	logger  logger.Logger
	storage Storage
}

func NewHandlers(logger logger.Logger, storage Storage) *Handlers {
	return &Handlers{logger: logger, storage: storage}
}

// Get returns the Value of a given Key
func (h *Handlers) Get(conn net.Conn, cmd []string) error {
	if len(cmd) != 2 {
		err := resp.NewError(fmt.Sprintf("ERR wrong number of arguments for 'GET' command"))
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

// Set sets a Key-Value pair in the storage
func (h *Handlers) Set(conn net.Conn, cmd []string) error {
	if len(cmd) != 3 {
		err := resp.NewError(fmt.Sprintf("ERR wrong number of arguments for 'SET' command"))
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

func (h *Handlers) Config(conn net.Conn, cmd []string) error {
	if len(cmd) == 1 {
		err := resp.NewError(fmt.Sprintf("ERR unable to return all the configuration"))

		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("unable to write all the configuration: %w", err)
		}
	}

	switch strings.Join(cmd, " ") {
	case "CONFIG GET appendonly":
		rsp := resp.NewArray([]string{"appendonly", "no"})
		if _, err := rsp.WriteTo(conn); err != nil {
			h.logger.Printf("unable to write: %v", err)
		}
	case "CONFIG GET save":
		rsp := resp.NewArray([]string{"save", "3600 1 300 100 60 10000"})
		if _, err := rsp.WriteTo(conn); err != nil {
			h.logger.Printf("unable to write: %v", err)
		}
	default:
		err := resp.NewError(fmt.Sprintf("ERR unsupported CONFIG command"))

		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("unsupported CONFIG command: %w", err)
		}
	}

	return nil

}

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
