package server

import (
	"context"
	"ddia/src/storage/aof"
	"errors"
	"io"
	"os"
)

// restoreAOF reads the AOF file and restores it into the server to keep the old state
func (s *Server) restoreAOF(ctx context.Context) error {
	aofPath := s.config.GetD("appenddirname", "./redis.aof")

	importAOF, err := aof.NewImportAppendOnlyFile(ctx, aofPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil // AOF does not exist. Nothing to import.
	} else if err != nil {
		return err
	}

	oldAOF := s.handlers.aof
	defer func() {
		s.handlers.aof = oldAOF // Resume the AOF configuration to the initial state after restoration
	}()
	s.handlers.aof = io.Discard // We don't want to record new records on the AOF when restoring the AOF!

	c := newClient(importAOF, s.options.dbs[0])
	c.authenticated = true // Pretend that we've successfully authenticated to the server

	return s.handleRequest(ctx, c)
}
