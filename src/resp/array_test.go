package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestBulkStr_ReadFrom(t *testing.T) {
	bulk := resp.Array{}

	input := strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	_, err := bulk.ReadFrom(input)
	require.NoError(t, err)

	require.Equal(t, "hello world", bulk.String())
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
			assert.Error(t, err)
			require.ErrorIs(t, err, resp.ErrParsingError)
			require.ErrorContains(t, err, tt.expectErrContains)
		})

	}
}

func TestBulkStr_StringAndBytes(t *testing.T) {
	bulk := resp.Array{}

	_, err := bulk.ReadFrom(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	require.NoError(t, err)

	require.Equal(t, "hello world", bulk.String())
	require.Equal(t, []byte(bulk.String()), bulk.Bytes())
}

func TestBulkStr_WriteTo(t *testing.T) {
	bulk := resp.Array{}

	text := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"

	_, err := bulk.ReadFrom(strings.NewReader(text))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	n, err := bulk.WriteTo(buf)
	require.NoError(t, err)
	require.NotZero(t, n)

	require.Equal(t, text, buf.String())
}
