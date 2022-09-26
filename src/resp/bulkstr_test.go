package resp_test

import (
	"ddia/src/resp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestBulkStr_ReadFrom(t *testing.T) {
	bulk := resp.BulkStr{}

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
			name:              "element of array is not string",
			input:             "*1\r\n:5\r\n",
			expectErrContains: "unknown operator",
		},
		{
			name:              "element is string but malformed",
			input:             "1\r\n$5\r\n",
			expectErrContains: "",
		},
	}

	bulk := resp.BulkStr{}

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
	bulk := resp.BulkStr{}

	_, err := bulk.ReadFrom(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	require.NoError(t, err)

	require.Equal(t, "hello world", bulk.String())
	require.Equal(t, []byte(bulk.String()), bulk.Bytes())
}
