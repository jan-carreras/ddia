package server_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"ddia/testing/log"
	"testing"
	"time"
)

func testServer() *server.Server {
	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger)
	return server.New(handlers, serverOptions()...)
}

func serverOptions() []server.Option {
	logger := log.ServerLogger()
	dbs := make([]server.Storage, 16)
	for i := 0; i < len(dbs); i++ {
		dbs[i] = storage.NewInMemory()
	}

	return []server.Option{
		server.WithLogger(logger),
		server.WithRandomPort(),
		server.WithDBs(dbs),
	}
}

func TestStart(t *testing.T) {
	s := testServer()

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	// TODO: This approach won't work if we're not adding support for reusing a connection
	// 	We would need to pass the s.Addr inside and manage the connection there (which makes total sense)
	cli := client.NewClient(log.ClientLogger(), s.Addr())

	rsp, err := cli.Set("hello", "world")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "OK"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}
}

func TestServer_Set(t *testing.T) {
	s := testServer()
	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	cli := client.NewClient(log.ClientLogger(), s.Addr())

	rsp, err := cli.Set("hello", "world")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "OK"; string(rsp) != want {
		t.Fatalf("invalid response: %q want %q", string(rsp), want)
	}

	v, err := cli.Get("hello")
	if err != nil {
		t.Fatalf("Set faield: %v, wanted no error", err)
	}
	if want := "world"; string(v) != want {
		t.Fatalf("invalid response: %q want %q", v, want)
	}
}

func TestServer_Ping(t *testing.T) {
	s := testServer()

	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("Start faield: %v, wanted no error", err)
	}

	cli := client.NewClient(log.ClientLogger(), s.Addr())

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
	s := testServer()

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
