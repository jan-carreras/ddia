package server_test

import (
	"context"
	"ddia/src/server"
	"ddia/testing/log"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestHandler_Auth(t *testing.T) {
	d := t.TempDir()
	f, err := os.Create(fmt.Sprintf("%s/redis.conf", d))
	if err != nil {
		t.Fatalf("unable to create temporary configuration file: %v", err)
	}

	_, err = f.WriteString("requirepass test-password-1234")
	if err != nil {
		t.Fatalf("unable to set password on config file: %v", err)
	}

	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger, io.Discard)

	options := serverOptions()
	options = append(options, server.WithConfigurationFile(f.Name()))

	s, err := server.New(handlers, options...)
	if err != nil {
		t.Fatalf("expecting server to be able to start without problems: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}

	t.Cleanup(func() { _ = s.Stop() })

	conn := testConn(t, s)

	rsp := req(t, conn, []string{"ping"})

	if want := "-NOAUTH Authentication required\r\n"; rsp != want {
		t.Fatalf("expecting to be deined: %q, want %q", rsp, want)
	}

	rsp = req(t, conn, []string{"auth", "invalid-password"})
	if want := "-WRONGPASS invalid username-password pair or user is disabled.\r\n"; rsp != want {
		t.Fatalf("expecting to be deined: %q, want %q", rsp, want)
	}

	rsp = req(t, conn, []string{"auth", "test-password-1234"})
	if want := "+OK\r\n"; rsp != want {
		t.Fatalf("expecting to be authenticated: %q, want %q", rsp, want)
	}

	rsp = req(t, conn, []string{"ping"})
	if want := "+PONG\r\n"; rsp != want {
		t.Fatalf("expecting to be authenticated: %q, want %q", rsp, want)
	}
}

func TestHandler_PingEcho(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("ping"), "PONG"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("ping hello world"), "hello world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("echo hello awesome world"), "hello awesome world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_Select(t *testing.T) {
	req := makeReq(t)

	// Database 0
	req("set hello world")

	// Database 1
	req("select 1")
	if rsp, want := req("get hello"), ""; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello there")
	if rsp, want := req("get hello"), "there"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	// Back to database 0 again
	req("select 0")

	if rsp, want := req("get hello"), "world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
