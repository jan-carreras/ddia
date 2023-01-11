package server_test

import (
	"fmt"
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
		want := fmt.Sprintf("+%d\r\n", n)
		if rsp != want {
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

	rsp, want = req("get hello"), "+world\r\n"
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

	rsp, want = req("ping hello world"), "+hello world\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("echo hello awesome world"), "+hello awesome world\r\n"
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
	rsp, want = req("get hello"), "+there\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	// Back to database 0 again
	req("select 0")
	rsp, want = req("get hello"), "+world\r\n"
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
