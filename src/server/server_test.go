package server_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"log"
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	logger = log.New(os.Stdout, "[client] ", 0)
	// TODO: This approach won't work if we're not adding support for reusing a connection
	// 	We would need to pass the s.Addr inside and manage the connection there (which makes total sense)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Set("hello", "world")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "OK"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}
}

func TestServer_Set(t *testing.T) {
	store := storage.NewInMemory()

	logger := log.New(os.Stdout, "[server] ", 0)
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	logger = log.New(os.Stdout, "[client] ", 0)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Set("hello", "world")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "OK"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}

	v, err := store.Get("hello")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "world"; v != want {
		t.Fatalf("invalid response: %q want %q", v, want)
	}
}

func TestServer_Ping(t *testing.T) {
	store := storage.NewInMemory()

	logger := log.New(os.Stdout, "[server] ", 0)
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	logger = log.New(os.Stdout, "[client] ", 0)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Ping("")
	if err != nil {
		t.Fatalf("Ping faield: %v, wanted no error", err)
	}

	if want := "PONG"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}

	rsp, err = cli.Ping("hello world")
	if err != nil {
		t.Fatalf("Ping faield: %v, wanted no error", err)
	}
	if want := "hello world"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}
}

func TestStart_GracefulShutdown(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger, store)

	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	time.Sleep(100 * time.Millisecond)
	err = s.Stop()
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}
}
