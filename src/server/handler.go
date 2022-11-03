package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"errors"
	"fmt"
	"net"
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

	// TODO: Understand which can of response should we return
	// 	It's unclear to me which data-type should I return:
	//  	Simple Strings or Bulk Strings?
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
