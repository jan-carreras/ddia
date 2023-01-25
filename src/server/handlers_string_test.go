package server_test

import (
	"bytes"
	"ddia/src/resp"
	"strconv"
	"testing"
)

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
