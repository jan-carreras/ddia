package server_test

import "testing"

func TestHandler_SetGetDel(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("set hello world"), "OK"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("get hello"), "world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("del hello"), "1"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("get hello"), ""; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_Exists(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("exists hello"), "0"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	if rsp, want := req("exists hello"), "1"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_RandomKey(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("randomkey"), "null"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	if rsp, want := req("randomkey"), "hello"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_Rename(t *testing.T) {
	req := makeReq(t)

	req("set hello world")
	req("rename hello new-hello")

	if rsp, want := req("get new-hello"), "world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	if rsp, want := req("exists hello"), "0"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
