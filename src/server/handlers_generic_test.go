package server_test

import "testing"

func TestHandler_SetGetDel(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("set hello world"), "OK"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("get hello"), "world"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("del hello"), "1"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("get hello"), ""
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_Exists(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("exists hello"), "0"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	rsp, want = req("exists hello"), "1"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_RandomKey(t *testing.T) {
	req := makeReq(t)

	rsp := req("randomkey")
	if want := "null"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set hello world")

	rsp = req("randomkey")
	if want := "hello"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_Rename(t *testing.T) {
	req := makeReq(t)

	req("set hello world")
	req("rename hello new-hello")
	rsp := req("get new-hello")

	if want := "world"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp = req("exists hello")
	if want := "0"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
