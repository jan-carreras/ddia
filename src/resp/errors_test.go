package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestError_ReadFrom_InvalidPrefix(t *testing.T) {
	s := resp.Error{}

	input := "(OK\r\n"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, 1, c) // Read the operation
}

func TestError_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.Error{}

	input := "-ERR unknown command 'helloworld'\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, len(input), c)
}

func TestError_ReadFrom(t *testing.T) {
	s := resp.Error{}

	input := "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)

	require.Equal(t, "WRONGTYPE Operation against a key holding the wrong kind of value", s.String())
}

func TestError_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "-" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestError_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize*2 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "-" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestError_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize*42 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "-" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestError_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.Error{}

	payloadSize := bufSize + (bufSize / 2) - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "-" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestError_WriteTo(t *testing.T) {
	input := "ERR Unexpected error"
	expected := "-ERR Unexpected error\r\n"
	s := resp.NewError(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	require.NoError(t, err)
	require.EqualValues(t, len(expected), c)

	require.Equal(t, expected, buf.String())
}
