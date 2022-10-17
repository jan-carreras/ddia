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
	time.Sleep(5 * time.Minute)
}

func startServer() error {
	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 6379, storage.NewInMemory())
	err := s.Start(context.Background())
	return err
}
