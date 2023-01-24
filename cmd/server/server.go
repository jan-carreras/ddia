// Package main is a command for the Redis server
package main

import (
	"context"
	"ddia/src/logger"
	"ddia/src/server"
	"ddia/src/server/config"
	"ddia/src/storage"
	aof2 "ddia/src/storage/aof"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() (err error) {
	ctx := context.Background()

	// Configuration
	cfg := config.NewEmpty()
	if len(os.Args) == 2 {
		cfg, err = config.New(os.Args[1]) // configPath
		if err != nil {
			return err
		}
	}

	l := getLogger(cfg)

	options, err := readOptions(cfg)
	if err != nil {
		return err
	}

	// Append Only File
	aof, err := getAOF(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() { _ = aof.Close() }()

	// HANDLERS
	handlers := server.NewHandlers(l, aof)

	// SERVER
	s, err := server.New(handlers, options...)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}

	// Start the server
	err = s.Start(context.Background())
	if err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	waitForGracefulShutdown(l, s)

	return nil
}

func getAOF(ctx context.Context, cfg config.Config) (aof io.WriteCloser, err error) {
	if cfg.GetD("appendonly", "no") != "yes" {
		return aof2.NewNoOpAOF(), nil // AOF is disabled
	}

	sync := aof2.NeverSync //nolint: ineffassign
	switch v := cfg.GetD("appendfsync", "always"); v {
	case "always":
		sync = aof2.AlwaysSync
	case "no":
		sync = aof2.NeverSync
	case "everysec":
		sync = aof2.EverySecondSync
	default:
		return nil, fmt.Errorf("appendfsync %q not supported", v)
	}

	f, err := os.OpenFile(
		cfg.GetD("appenddirname", "./redis.aof"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return nil, err
	}

	return aof2.NewAppendOnlyFile(ctx, f, sync), nil
}

func getLogger(_ config.Config) logger.Logger {
	return log.New(os.Stdout, "[server] ", 0)
}

func readOptions(cfg config.Config) ([]server.Option, error) {
	var options []server.Option

	// Configuration
	if len(os.Args) == 2 {
		configPath := os.Args[1]
		options = append(options, server.WithConfigurationFile(configPath))
	}

	// Logger
	options = append(options, server.WithLogger(getLogger(cfg)))

	// DBs
	databases, err := cfg.Integer("databases", 16)
	if err != nil {
		return nil, err
	}
	dbs := make([]server.Storage, databases)
	for i := 0; i < len(dbs); i++ {
		dbs[i] = storage.NewInMemory()
	}
	options = append(options, server.WithDBs(dbs))

	// Port
	port, err := cfg.Integer("port", 6379)
	if err != nil {
		return nil, err
	}
	options = append(options, server.WithPort(port))

	return options, nil
}

func waitForGracefulShutdown(logger logger.Logger, srv *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	if err := srv.Stop(); err != nil {
		logger.Println("Server forced to shutdown: %v", err)
	}

	logger.Printf("Server stopped")
}
