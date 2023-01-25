package server_test

import "testing"

func TestHandler_DBSize(t *testing.T) {
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

func TestHandler_FlushAll(t *testing.T) {
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

func TestHandler_FlushDB(t *testing.T) {
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
