package client_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"ddia/testing/log"
	"testing"
)

func TestClient_Set(t *testing.T) {
	logger := log.ServerLogger()
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger)
	s, err := server.New(
		handlers,
		server.WithRandomPort(),
		server.WithLogger(logger),
		server.WithDBs([]server.Storage{store}),
	)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("start application: %v", err)
	}

	c := client.NewClient(log.ClientLogger(), s.Addr())

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
	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger)
	s, err := server.New(
		handlers,
		server.WithRandomPort(),
		server.WithLogger(logger),
		server.WithDBs([]server.Storage{store}),
	)
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	c := client.NewClient(log.ClientLogger(), s.Addr())

	k, v := "hello", "world"
	err = store.Set(k, v)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	rsp, err := c.Get(k)
	if err != nil {
		t.Fatalf("expecting no error, got: %q", err)
	}

	if want := "world"; want != string(rsp) {
		t.Fatalf("invalid response: %q, expecting %q", string(rsp), want)
	}
}

func TestClient_Ping(t *testing.T) {
	store := storage.NewInMemory()
	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger)
	s, err := server.New(
		handlers,
		server.WithRandomPort(),
		server.WithLogger(logger),
		server.WithDBs([]server.Storage{store}),
	)

	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	c := client.NewClient(log.ClientLogger(), s.Addr())

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
