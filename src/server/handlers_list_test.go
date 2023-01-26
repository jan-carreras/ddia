package server_test

import "testing"

func TestListOperations(t *testing.T) {
	req := makeReq(t)

	c := req("llen mylist")
	if want := ":0\r\n"; c != want {
		t.Fatalf("expecting empty list: %q, want %q", c, want)
	}

	// Lpush appends on the start of the list
	req("lpush mylist three") // -> three
	req("lpush mylist two")   // -> two three
	req("lpush mylist one")   // -> one two three

	// rpush appends at the end of the list
	req("rpush mylist four") // one two three four          <-
	req("rpush mylist five") // one two three four five     <-
	req("rpush mylist six")  // one two three four five six <-

	c = req("llen mylist")
	if want := ":6\r\n"; c != want {
		t.Fatalf("unexpected list length: %q, want %q", c, want)
	}

	list := req("lrange mylist 0 5")
	want := "*6\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n$4\r\nfour\r\n$4\r\nfive\r\n$3\r\nsix\r\n"
	if list != want {
		t.Fatalf("unexpected list: %q, want %q", list, want)
	}

	req("rpush mylist ***")
	req("rpush mylist six")
	req("rpush mylist six") // one two three four five six six six <-

	c = req("lrem mylist -2 six") // delete two 6 from the right
	if want := ":2\r\n"; c != want {
		t.Fatalf("unexpected items deleted: %q, want %q", c, want)
	}

	list = req("lrange mylist 0 -1")
	want = "*7\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n$4\r\nfour\r\n$4\r\nfive\r\n$3\r\nsix\r\n$3\r\n***\r\n"
	if list != want {
		t.Fatalf("unexpected list: %q, want %q", list, want)
	}

	end := req("rpop mylist")
	want = "$3\r\n***\r\n"
	if end != want {
		t.Fatalf("expecting empty list: %q, want %q", end, want)
	}

	start := req("lpop mylist")
	want = "$3\r\none\r\n"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	start = req("lindex mylist 0")
	want = "$3\r\ntwo\r\n"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	end = req("lindex mylist -1")
	want = "$3\r\nsix\r\n"
	if end != want {
		t.Fatalf("unexpected value returned: %q, want %q", end, want)
	}

	rsp := req("lset mylist 0 newvalue")
	want = "$2\r\nOK\r\n"
	if rsp != want {
		t.Fatalf("unexpected value returned: %q, want %q", rsp, want)
	}

	start = req("lindex mylist 0")
	want = "$8\r\nnewvalue\r\n"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	rsp = req("lset mylist -1 newvalueEnd")
	want = "$2\r\nOK\r\n"
	if rsp != want {
		t.Fatalf("unexpected value returned: %q, want %q", rsp, want)
	}

	start = req("lindex mylist -1")
	want = "$11\r\nnewvalueEnd\r\n"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}
}
