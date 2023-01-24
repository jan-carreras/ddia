package server_test

import (
	"bytes"
	"context"
	"ddia/src/resp"
	"ddia/src/server"
	"ddia/testing/log"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestServer_UnknownCommand(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("foo"), "-ERR unknown command 'foo'\r\n"
	if rsp != want {
		t.Fatalf("missmatch: %q, want %q", rsp, want)
	}
}

func TestServer_IncrDecrOperators(t *testing.T) {
	req := makeReq(t)

	assert := func(n int) {
		rsp := req("get key")
		s := resp.NewStr(strconv.Itoa(n))
		want := &bytes.Buffer{}
		_, _ = s.WriteTo(want)
		if rsp != want.String() {
			t.Fatalf("missmatch: %q, want %q", rsp, want)
		}
	}

	req("incr key")
	assert(1)

	req("decr key")
	assert(0)

	req("incrby key 11")
	assert(11)

	req("decrby key 10")
	assert(1)
}

func TestServer_SetGetDel(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("set hello world"), "+OK\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("get hello"), "$5\r\nworld\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("del hello"), ":1\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("get hello"), "+\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_PingEcho(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("ping"), "+PONG\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("ping hello world"), "$11\r\nhello world\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("echo hello awesome world"), "$19\r\nhello awesome world\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_Select(t *testing.T) {
	req := makeReq(t)

	// Database 0
	req("set hello world")

	// Database 1
	req("select 1")
	rsp, want := req("get hello"), "+\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
	req("set hello there")
	rsp, want = req("get hello"), "$5\r\nthere\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	// Back to database 0 again
	req("select 0")
	rsp, want = req("get hello"), "$5\r\nworld\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_DBSize(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("dbsize"), ":0\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set one 1")
	req("set two 2")
	req("set three 3")

	rsp, want = req("dbsize"), ":3\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestStart_GracefulShutdown(t *testing.T) {
	s := testServer(t)

	time.Sleep(10 * time.Millisecond)

	err := s.Stop()
	if err != nil {
		t.Fatalf("Stop faield: %v, wanted no error", err)
	}
}

func TestServer_Auth(t *testing.T) {
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

func TestServer_FlushDB(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("flushdb"), "+OK\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	t.Run("remove actual data", func(t *testing.T) {
		req("set hello world")
		req("set answer-to-everything 42")

		rsp, want := req("dbsize"), ":2\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		rsp, want = req("flushdb"), "+OK\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		// Database is empty after flush
		rsp, want = req("dbsize"), ":0\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}
	})
}

func TestServer_FlushAll(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("flushall"), "+OK\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	t.Run("remove actual data", func(t *testing.T) {
		req("set hello world")
		req("set answer-to-everything 42")
		req("select 1")
		req("set hello world")

		rsp, want := req("dbsize"), ":1\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		rsp, want = req("flushall"), "+OK\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		// Database is empty after flush
		rsp, want = req("dbsize"), ":0\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		req("select 0")
		rsp, want = req("dbsize"), ":0\r\n"
		if rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}
	})
}

func TestServer_Exists(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("exists hello"), ":0\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	rsp, want = req("exists hello"), ":1\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_MGet(t *testing.T) {
	req := makeReq(t)

	req("set one one")
	req("set two two")

	rsp := req("mget one two")

	want := "*2\r\n$3\r\none\r\n$3\r\ntwo\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_SetNX(t *testing.T) {
	req := makeReq(t)

	rsp := req("setnx hello world")
	if want := ":1\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp = req("setnx hello universe")
	if want := ":0\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp = req("get hello")
	if want := "$5\r\nworld\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_RandomKey(t *testing.T) {
	req := makeReq(t)

	rsp := req("randomkey")
	if want := "$-1\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	rsp = req("randomkey")
	if want := "$5\r\nhello\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestServer_Rename(t *testing.T) {
	req := makeReq(t)

	req("set hello world")
	req("rename hello new-hello")
	rsp := req("get new-hello")

	if want := "$5\r\nworld\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp = req("exists hello")
	if want := ":0\r\n"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
