package server_test

import "testing"

func TestListOperations(t *testing.T) {
	req := makeReq(t)

	if have, want := req("llen mylist"), "0"; have != want {
		t.Fatalf("expecting empty list: %q, want %q", have, want)
	}

	// Lpush appends on the start of the list
	req("lpush mylist three") // -> three
	req("lpush mylist two")   // -> two three
	req("lpush mylist one")   // -> one two three

	// rpush appends at the end of the list
	req("rpush mylist four") // one two three four          <-
	req("rpush mylist five") // one two three four five     <-
	req("rpush mylist six")  // one two three four five six <-

	if have, want := req("llen mylist"), "6"; have != want {
		t.Fatalf("unexpected list length: %q, want %q", have, want)
	}

	if have, want := req("lrange mylist 0 5"), "one two three four five six"; have != want {
		t.Fatalf("unexpected list: %q, want %q", have, want)
	}

	req("rpush mylist ***")
	req("rpush mylist six")
	req("rpush mylist six") // one two three four five six *** six six <-

	// delete two 6 from the right
	if have, want := req("lrem mylist -2 six"), "2"; have != want {
		t.Fatalf("unexpected items deleted: %q, want %q", have, want)
	}

	if have, want := req("lrange mylist 0 -1"), "one two three four five six ***"; have != want {
		t.Fatalf("unexpected list: %q, want %q", have, want)
	}

	if have, want := req("rpop mylist"), "***"; have != want {
		t.Fatalf("expecting empty list: %q, want %q", have, want)
	}

	if have, want := req("lpop mylist"), "one"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lindex mylist 0"), "two"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lindex mylist -1"), "six"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lset mylist 0 newvalue"), "OK"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lindex mylist 0"), "newvalue"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lset mylist -1 newvalueEnd"), "OK"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("lindex mylist -1"), "newvalueEnd"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	if have, want := req("ltrim mylist 0 1"), "OK"; have != want {
		t.Fatalf("unexpected value returned: %q, want %q", have, want)
	}

	// create: one two three four
	req("rpush mylist2 one")
	req("rpush mylist2 two")
	req("rpush mylist2 three")
	req("rpush mylist2 four")

	if have, want := req("lrange mylist2 1 2"), "two three"; have != want {
		t.Fatalf("unexpected list: %q, want %q", have, want)
	}
}
