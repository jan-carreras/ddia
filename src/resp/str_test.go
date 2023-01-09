package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"errors"
	"strings"
	"testing"
)

func TestStr_ReadFrom(t *testing.T) {
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader("$5\r\nhello\r\n"))
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	if want := "hello"; s.String() != want {
		t.Fatalf("invalid response: %q, want %q", s.String(), want)
	}
}

func TestStr_ReadFrom_EmptyString(t *testing.T) {
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader("$0\r\n\r\n"))
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	if want := ""; s.String() != want {
		t.Fatalf("invalid response: %q want %q", s.String(), want)
	}
}

func TestStr_ReadFrom_Errors(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedErrContains string
	}{
		{
			name:                "missing operation type",
			input:               "\r\nhello\r\n",
			expectedErrContains: "unknown operation",
		},
		{
			name:                "invalid operation type",
			input:               "?\r\nhello\r\n",
			expectedErrContains: "unknown operation",
		},
		{
			name:                "invalid length",
			input:               "$\r\nhello\r\n",
			expectedErrContains: "readLength",
		},
		{
			name:                "length and string mismatch: string too short",
			input:               "$10\r\nhello\r\n",
			expectedErrContains: "insufficient data read",
		},
		{
			name:                "length and string mismatch: string too long",
			input:               "$5\r\nhello world\r\n",
			expectedErrContains: "unexpected character",
		},
	}

	s := resp.Str{}
	for _, tt := range tests {
		_, err := s.ReadFrom(strings.NewReader(tt.input))
		if !errors.Is(err, resp.ErrParsingError) {
			t.Fatalf("invalid error type: %v, want %v", err, resp.ErrParsingError)
		}
		if !strings.Contains(err.Error(), tt.expectedErrContains) {
			t.Fatalf("invalid error content: %q want %q", err.Error(), tt.expectedErrContains)
		}
	}
}

func TestStr_WriteTo(t *testing.T) {
	original := "$5\r\nhello\r\n"
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader(original))
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	if want := original; buf.String() != want {
		t.Fatalf("invalid response: %q want %q", buf.String(), want)
	}
}

func TestStr_EmptyString(t *testing.T) {
	original := "$0\r\n\r\n"
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader(original))
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	if err != nil {
		t.Fatalf("error unexpected: %v. want no error", err)
	}

	if want := original; buf.String() != want {
		t.Fatalf("invalid response: %q want %q", buf.String(), want)
	}
}
