package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"strings"
	"testing"
)

func TestInteger_ReadFrom_InvalidPrefix(t *testing.T) {
	s := resp.Integer{}

	input := "(1234\r\n"
	c, err := s.ReadFrom(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expecting error, got %v", err)
	}
	if want := 1; int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
}

func TestInteger_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.Integer{}

	input := ":1234\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	if err == nil {
		t.Fatalf("expecting error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
}

func TestInteger_ReadFrom(t *testing.T) {
	s := resp.Integer{}

	input := ":1234\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}

	if want := "1234"; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}

func TestInteger_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("5", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}

func TestInteger_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize*2 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}

func TestInteger_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize*42 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("6", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}

func TestInteger_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize + (bufSize / 2) - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("4", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(input); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}
	if want := payload; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}

func TestInteger_WriteTo(t *testing.T) {
	input := 42
	expected := ":42\r\n"
	s := resp.NewInteger(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	if err != nil {
		t.Fatalf("expecting no error, got %v", err)
	}
	if want := len(expected); int(c) != want {
		t.Fatalf("invalid number of characters read %d, want %d", c, want)
	}

	if want := "42"; s.String() != want {
		t.Fatalf("invalid String parsed: %q want %q", s.String(), want)
	}
}
