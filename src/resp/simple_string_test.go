package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"strings"
	"testing"
)

const bufSize = 512

func TestSimpleString_ReadFrom_InvalidPrefix(t *testing.T) {
	s := resp.SimpleString{}

	input := "(OK\r\n"
	c, err := s.ReadFrom(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expecting error, got no error instead")
	}

	if want := 1; int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
}

func TestSimpleString_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.SimpleString{}

	input := "+OK\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expecting error, got no error instead")
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
}

func TestSimpleString_ReadFrom(t *testing.T) {
	s := resp.SimpleString{}

	input := "+OK\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}

	if want := "OK"; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestSimpleString_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestSimpleString_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize*2 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestSimpleString_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize*42 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}

	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestSimpleString_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize + (bufSize / 2) - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}

	if want := len(input); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestSimpleString_WriteTo(t *testing.T) {
	input := "OK"
	expected := "+OK\r\n"
	s := resp.NewSimpleString(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}

	if want := len(expected); int(c) != want {
		t.Fatalf("invalid readCount: %d, want %d", c, want)
	}
	if want := expected; buf.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}
