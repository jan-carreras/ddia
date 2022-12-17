package main

import (
	"context"
	"ddia/src/server"
	"ddia/src/storage"
	"log"
	"os"
	"time"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Minute) // TODO: Actively wait
}

func startServer() error {
	logger := log.New(os.Stdout, "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 6379, handlers)
	err := s.Start(context.Background())
	return err
}
