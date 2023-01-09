package main

import (
	"context"
	"ddia/src/logger"
	"ddia/src/server"
	"ddia/src/storage"
	"fmt"
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
	log := log.New(os.Stdout, "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(log, store)
	s := server.New(handlers)

	err := s.Start(context.Background())
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	waitForGracefulShutdown(log, s)

	return nil
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
