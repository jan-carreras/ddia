package server_test

import (
	"testing"
)

func TestHandler_IncrDecrOperators(t *testing.T) {
	req := makeReq(t)

	assert := func(want string) {
		rsp := req("get key")
		if rsp != want {
			t.Fatalf("missmatch: %q, want %q", rsp, want)
		}
	}

	req("incr key")
	assert("1")

	req("decr key")
	assert("0")

	req("incrby key 11")
	assert("11")

	req("decrby key 10")
	assert("1")
}

func TestHandler_MGet(t *testing.T) {
	req := makeReq(t)

	req("set one one")
	req("set two two")

	rsp := req("mget one two")

	want := "one two"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_SetNX(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("setnx hello world"), "1"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("setnx hello universe"), "0"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("get hello"), "world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
