package server_test

import (
	"bufio"
	"context"
	"ddia/src/resp"
	"ddia/src/server"
	"ddia/src/storage"
	"ddia/testing/log"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
)

func parse(t testing.TB, response string) fmt.Stringer {
	dt, err := resp.Decode(strings.NewReader(response))
	if err != nil {
		t.Fatal(err)
	}

	return dt
}

func makeReq(t testing.TB) func(string) string {
	s := testServer(t)
	conn := testConn(t, s)

	return func(args string) string {
		response := req(t, conn, strings.Split(args, " "))
		return parse(t, response).String()
	}
}

func req(t testing.TB, conn net.Conn, req []string) string {
	reader := bufio.NewReader(conn)

	r := resp.NewArray(req)
	_, err := r.WriteTo(conn)
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
	}

	buf := make([]byte, 1024*10) // This is going to byte my ass, for sure

	n, err := reader.Read(buf)
	if errors.Is(err, io.EOF) {
	} else if err != nil {
		t.Fatalf("not expecting error: %q", err.Error())
	}

	return string(buf[:n])
}

func testConn(t testing.TB, s *server.Server) net.Conn {
	conn, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
	}
	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

func testServer(t testing.TB) *server.Server {
	t.Helper()

	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger, io.Discard)

	s, err := server.New(handlers, serverOptions()...)
	if err != nil {
		t.Fatalf("expecting server to be able to start without problems: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
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
