package client_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"io"
	"log"
	"testing"
)

func loggerOutput() io.Writer {
	return io.Discard
}

func TestClient_Set(t *testing.T) {
	logger := log.New(loggerOutput(), "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("start application: %v", err)
	}

	logger = log.New(loggerOutput(), "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	rsp, err := c.Set("hello", "world")
	if err != nil {
		t.Fatalf("c.Set: %v", err)
	}

	if want := "OK"; string(rsp) != want {
		t.Fatalf("response: %q, expecting %q", string(rsp), want)
	}

	rsp, err = c.Set("chao", "universe")
	if err != nil {
		t.Fatalf("c.Set: %v", err)
	}

	if want := "OK"; string(rsp) != want {
		t.Fatalf("response: %q, expecting %q", string(rsp), want)
	}
}

func TestClient_Get(t *testing.T) {
	store := storage.NewInMemory()
	logger := log.New(loggerOutput(), "[server] ", 0)
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	logger = log.New(loggerOutput(), "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	k, v := "hello", "world"
	err = store.Set(k, v)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	rsp, err := c.Get(k)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if want := "world"; want != string(rsp) {
		t.Fatalf("invalid response: %q, expecting %q", string(rsp), want)
	}
}

func TestClient_Ping(t *testing.T) {
	store := storage.NewInMemory()
	logger := log.New(loggerOutput(), "[server] ", 0)
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	logger = log.New(loggerOutput(), "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	rsp, err := c.Ping("")
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}

	if want := "PONG"; string(rsp) != want {
		t.Fatalf("invalid response: %q, expecting %q", string(rsp), want)
	}

	rsp, err = c.Ping("hello world")
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}

	if want := "hello world"; want != string(rsp) {
		t.Fatalf("invalid response: %q, expecting %q", string(rsp), want)
	}
}
