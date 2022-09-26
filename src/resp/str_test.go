package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestStr_ReadFrom(t *testing.T) {
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader("$5\r\nhello\r\n"))
	require.NoError(t, err)

	require.Equal(t, "hello", s.String())
}

func TestStr_ReadFrom_EmptyString(t *testing.T) {
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader("$0\r\n\r\n"))
	require.NoError(t, err)

	require.Equal(t, "", s.String())
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
		require.ErrorIs(t, err, resp.ErrParsingError)
		require.ErrorContains(t, err, tt.expectedErrContains)
	}

}

func TestStr_WriteTo(t *testing.T) {
	original := "$5\r\nhello\r\n"
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader(original))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	require.NoError(t, err)

	require.Equal(t, original, buf.String())
}

func TestStr_EmptyString(t *testing.T) {
	original := "$0\r\n\r\n"
	s := resp.Str{}
	_, err := s.ReadFrom(strings.NewReader(original))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	_, err = s.WriteTo(buf)
	require.NoError(t, err)

	require.Equal(t, original, buf.String())
}
