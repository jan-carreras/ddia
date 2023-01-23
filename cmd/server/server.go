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

func startServer() error {
	// Configuration
	cfg := config.NewEmpty()
	if len(os.Args) == 2 {
		configPath := os.Args[1]

		var err error
		cfg, err = config.New(configPath)
		if err != nil {
			return err
		}
	}

	l, err := getLogger(cfg)
	if err != nil {
		return err
	}

	options, err := readOptions(cfg)
	if err != nil {
		return err
	}

	// Append Only File
	aof := io.Discard
	if cfg.GetD("appendonly", "no") == "yes" {
		sync := aof2.NeverSync //nolint: ineffassign
		switch v := cfg.GetD("appendfsync", "always"); v {
		case "always":
			sync = aof2.AlwaysSync
		case "no":
			sync = aof2.NeverSync
		case "everysec":
			sync = aof2.EverySecondSync
		default:
			return fmt.Errorf("appendfsync %q not supported", v)
		}

		appenddirname := "./redis.aof"
		if path, ok := cfg.Get("appenddirname"); ok {
			appenddirname = path
		}
		f, err := os.OpenFile(appenddirname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		aof = aof2.NewAppendOnlyFile(f, sync)
	}

	handlers := server.NewHandlers(l, aof)

	//
	s, err := server.New(handlers, options...)

	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	waitForGracefulShutdown(l, s)

	return nil
}

func getLogger(_ config.Config) (logger.Logger, error) {
	l := log.New(os.Stdout, "[server] ", 0)
	return l, nil
}

func readOptions(cfg config.Config) ([]server.Option, error) {
	var options []server.Option

	// Configuration
	if len(os.Args) == 2 {
		configPath := os.Args[1]
		options = append(options, server.WithConfigurationFile(configPath))
	}

	// Logger
	l := log.New(os.Stdout, "[server] ", 0)
	options = append(options, server.WithLogger(l))

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
