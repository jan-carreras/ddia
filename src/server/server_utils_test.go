package server_test

import (
	"context"
	"ddia/src/resp"
	"ddia/src/server"
	"ddia/src/storage"
	"ddia/testing/log"
	"net"
	"strings"
	"testing"
)

func makeReq(t *testing.T) func(string) string {
	s := testServer(t)
	conn := testConn(t, s)

	return func(args string) string {
		return req(t, conn, strings.Split(args, " "))
	}
}

func req(t *testing.T, conn net.Conn, req []string) string {
	t.Helper()

	r := resp.NewArray(req)
	_, err := r.WriteTo(conn)
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}

	buf := make([]byte, 1024) // This is going to byte my ass, for sure

	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("not expecing error: %q", err.Error())
	}

	return string(buf[:n])
}

func testConn(t *testing.T, s *server.Server) net.Conn {
	conn, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}
	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

func testServer(t *testing.T) *server.Server {
	t.Helper()

	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger)

	s, err := server.New(handlers, serverOptions()...)
	if err != nil {
		t.Fatalf("expecting server to be able to start without problems: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}

	t.Cleanup(func() { _ = s.Stop() })
	return s
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