package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"errors"
	"strings"
	"testing"
)

func TestBulkStr_ReadFrom(t *testing.T) {
	bulk := resp.Array{}

	input := strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	_, err := bulk.ReadFrom(input)
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}

	if want := "hello world"; bulk.String() != want {
		t.Fatalf("invalid response: %q, want %q", bulk.String(), want)
	}
}

func TestBulkStr_ReadFrom_Errors(t *testing.T) {

	tests := []struct {
		name              string
		input             string
		expectErrContains string
	}{
		{
			name:              "missing array length",
			input:             "*\r\n$5\r\nhello\r\n",
			expectErrContains: "readLength",
		},
		{
			name:              "empty input",
			input:             "*",
			expectErrContains: "readLength",
		},
		{
			name:              "missing elements",
			input:             "*10\r\n",
			expectErrContains: "unable to read operator",
		},
		{
			name:              "element is string but malformed",
			input:             "*1\r\n$5\r\n",
			expectErrContains: "",
		},
	}

	bulk := resp.Array{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bulk.ReadFrom(strings.NewReader(tt.input))
			if err == nil {
				t.Errorf("expecting error: %v", err)
			}

			if !errors.Is(err, resp.ErrParsingError) {
				t.Fatalf("invalid error returned: %v, expecting %v", err, resp.ErrParsingError)
			}

			if !strings.Contains(err.Error(), tt.expectErrContains) {
				t.Fatalf("invalid error message: %v, expected to contain %q", err, tt.expectErrContains)
			}
		})

	}
}

func TestBulkStr_StringAndBytes(t *testing.T) {
	bulk := resp.Array{}

	_, err := bulk.ReadFrom(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}

	if want := "hello world"; bulk.String() != want {
		t.Fatalf("invalid response: %q, expecting %q", bulk.String(), want)
	}

	if !bytes.Equal([]byte(bulk.String()), bulk.Bytes()) {
		t.Fatalf("bytes and string response must be the same: %q want %q", string(bulk.Bytes()), bulk.String())

	}
}

func TestBulkStr_WriteTo(t *testing.T) {
	bulk := resp.Array{}

	text := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"

	_, err := bulk.ReadFrom(strings.NewReader(text))
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}

	buf := &bytes.Buffer{}
	n, err := bulk.WriteTo(buf)
	if err != nil {
		t.Fatalf("WriteTo: %v", err)
	}

	if n == 0 {
		t.Fatalf("invalid n: %d, want non-zero instead", n)
	}

	if buf.String() != text {
		t.Fatalf("invalid text: %q, want %q", buf.String(), text)
	}
}
