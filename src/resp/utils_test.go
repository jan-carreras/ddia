package resp

import (
	"strings"
	"testing"
)

func TestReadLength(t *testing.T) {
	length, err := readLength(strings.NewReader("12\r\n"))
	if err != nil {
		t.Fatalf("error not expected: %q", err)
	}

	if want := 12; length != want {
		t.Fatalf("unexpected length: %d, want %d", length, want)
	}

	// Null string
	length, err = readLength(strings.NewReader("-1\r\n"))
	if err != nil {
		t.Fatalf("error not expected: %q", err)
	}

	if want := -1; length != want {
		t.Fatalf("unexpected length: %d, want %d", length, want)
	}
}
