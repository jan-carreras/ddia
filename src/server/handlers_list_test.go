package server_test

import "testing"

func TestListOperations(t *testing.T) {
	req := makeReq(t)

	c := req("llen mylist")
	if want := "0"; c != want {
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
	if want := "6"; c != want {
		t.Fatalf("unexpected list length: %q, want %q", c, want)
	}

	list := req("lrange mylist 0 5")
	want := "one two three four five six"
	if list != want {
		t.Fatalf("unexpected list: %q, want %q", list, want)
	}

	req("rpush mylist ***")
	req("rpush mylist six")
	req("rpush mylist six") // one two three four five six six six <-

	c = req("lrem mylist -2 six") // delete two 6 from the right
	if want := "2"; c != want {
		t.Fatalf("unexpected items deleted: %q, want %q", c, want)
	}

	list = req("lrange mylist 0 -1")
	want = "one two three four five six ***"
	if list != want {
		t.Fatalf("unexpected list: %q, want %q", list, want)
	}

	end := req("rpop mylist")
	want = "***"
	if end != want {
		t.Fatalf("expecting empty list: %q, want %q", end, want)
	}

	start := req("lpop mylist")
	want = "one"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	start = req("lindex mylist 0")
	want = "two"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	end = req("lindex mylist -1")
	want = "six"
	if end != want {
		t.Fatalf("unexpected value returned: %q, want %q", end, want)
	}

	rsp := req("lset mylist 0 newvalue")
	want = "OK"
	if rsp != want {
		t.Fatalf("unexpected value returned: %q, want %q", rsp, want)
	}

	start = req("lindex mylist 0")
	want = "newvalue"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}

	rsp = req("lset mylist -1 newvalueEnd")
	want = "OK"
	if rsp != want {
		t.Fatalf("unexpected value returned: %q, want %q", rsp, want)
	}

	start = req("lindex mylist -1")
	want = "newvalueEnd"
	if start != want {
		t.Fatalf("unexpected value returned: %q, want %q", start, want)
	}
}
