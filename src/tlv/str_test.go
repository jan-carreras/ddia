package tlv_test

import (
	"ddia/src/tlv"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestStr_ReadFrom(t *testing.T) {
	s := tlv.Str{}
	_, err := s.ReadFrom(strings.NewReader("5\r\nhello\r\n"))
	require.NoError(t, err)

	require.Equal(t, "hello", s.String())
}

func TestStr_ReadFrom_Errors(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedErrContains string
	}{
		{
			name:                "invalid length",
			input:               "\r\nhello\r\n",
			expectedErrContains: "readLength",
		},
		{
			name:                "length and string mismatch: string too short",
			input:               "10\r\nhello\r\n",
			expectedErrContains: "insufficient data read",
		},
		{
			name:                "length and string mismatch: string too long",
			input:               "5\r\nhello world\r\n",
			expectedErrContains: "unexpected character",
		},
	}

	s := tlv.Str{}
	for _, tt := range tests {
		_, err := s.ReadFrom(strings.NewReader(tt.input))
		require.ErrorIs(t, err, tlv.ErrParsingError)
		require.ErrorContains(t, err, tt.expectedErrContains)
	}

}
