package server_test

import "testing"

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
