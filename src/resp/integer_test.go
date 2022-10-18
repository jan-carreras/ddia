package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestInteger_ReadFrom_InvalidPrefix(t *testing.T) {
	s := resp.Integer{}

	input := "(1234\r\n"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, 1, c) // Read the operation
}

func TestInteger_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.Integer{}

	input := ":1234\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, len(input), c)
}

func TestInteger_ReadFrom(t *testing.T) {
	s := resp.Integer{}

	input := ":1234\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)

	require.Equal(t, "1234", s.String())
}

func TestInteger_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("5", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestInteger_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize*2 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestInteger_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize*42 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("6", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestInteger_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.Integer{}

	payloadSize := bufSize + (bufSize / 2) - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("4", payloadSize)
	input := ":" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestInteger_WriteTo(t *testing.T) {
	input := "42"
	expected := ":42\r\n"
	s := resp.NewInteger(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	require.NoError(t, err)
	require.EqualValues(t, len(expected), c)

	require.Equal(t, expected, buf.String())
}
