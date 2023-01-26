package server_test

import "testing"

func TestHandler_DBSize(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("dbsize"), "0"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	req("set one 1")
	req("set two 2")
	req("set three 3")

	if rsp, want := req("dbsize"), "3"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}

func TestHandler_FlushAll(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("flushall"), "OK"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	t.Run("remove actual data", func(t *testing.T) {
		req("set hello world")
		req("set answer-to-everything 42")
		req("select 1")
		req("set hello world")

		if rsp, want := req("dbsize"), "1"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		if rsp, want := req("flushall"), "OK"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		// Database is empty after flush
		if rsp, want := req("dbsize"), "0"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		req("select 0")
		if rsp, want := req("dbsize"), "0"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}
	})
}

func TestHandler_FlushDB(t *testing.T) {
	req := makeReq(t)

	if rsp, want := req("flushdb"), "OK"; rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	t.Run("remove actual data", func(t *testing.T) {
		req("set hello world")
		req("set answer-to-everything 42")

		if rsp, want := req("dbsize"), "2"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		if rsp, want := req("flushdb"), "OK"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}

		// Database is empty after flush
		if rsp, want := req("dbsize"), "0"; rsp != want {
			t.Fatalf("invalid response: %q want %q", rsp, want)
		}
	})
}
