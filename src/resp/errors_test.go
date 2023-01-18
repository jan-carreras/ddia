package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"strings"
	"testing"
)

func TestError_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.Error{}

	input := "ERR unknown command 'helloworld'\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expecting error, got nil instead")
	}

	if len(input) != int(c) {
		t.Fatalf("invalid readCount %d, want %d", c, len(input))
	}
}

func TestError_ReadFrom(t *testing.T) {
	s := resp.Error{}

	input := "WRONGTYPE Operation against a key holding the wrong kind of value\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}

	if len(input) != int(c) {
		t.Fatalf("invalid readCount %d, want %d", c, len(input))
	}

	want := "WRONGTYPE Operation against a key holding the wrong kind of value"
	if s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestError_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize - 2 // -2 because \r\n
	payload := strings.Repeat("h", payloadSize)
	input := payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}
	if int(c) != len(input) {
		t.Fatalf("different lengths: %d want %d", c, len(input))
	}
	if s.String() != payload {
		t.Fatalf("Strings() not equal %q, want %q", s.String(), payload)
	}
}

func TestError_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize*2 - 2 // -2 because \r\n
	payload := strings.Repeat("h", payloadSize)
	input := payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}
	if int(c) != len(input) {
		t.Fatalf("different lengths: %d want %d", c, len(input))
	}
	if s.String() != payload {
		t.Fatalf("Strings() not equal %q, want %q", s.String(), payload)
	}
}

func TestError_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize*42 - 2 // -2 because \r\n
	payload := strings.Repeat("h", payloadSize)
	input := payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}
	if int(c) != len(input) {
		t.Fatalf("different lengths: %d want %d", c, len(input))
	}
	if s.String() != payload {
		t.Fatalf("Strings() not equal %q, want %q", s.String(), payload)
	}
}

func TestError_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize + (bufSize / 2) - 2 // -2 because \r\n
	payload := strings.Repeat("h", payloadSize)
	input := payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}
	if int(c) != len(input) {
		t.Fatalf("different lengths: %d want %d", c, len(input))
	}
	if s.String() != payload {
		t.Fatalf("Strings() not equal %q, want %q", s.String(), payload)
	}
}

func TestError_WriteTo(t *testing.T) {
	input := "ERR Unexpected error"
	expected := "-ERR Unexpected error\r\n"
	s := resp.NewError(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	if err != nil {
		t.Fatalf("expecing not error, got %v", err)
	}
	if int(c) != len(expected) {
		t.Fatalf("different lengths: %d want %d", c, len(expected))
	}
	if buf.String() != expected {
		t.Fatalf("Strings() not equal %q, want %q", s.String(), expected)
	}
}
