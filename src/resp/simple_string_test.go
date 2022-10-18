package resp_test

import (
	"bytes"
	"ddia/src/resp"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const bufSize = 512

func TestSimpleString_ReadFrom_InvalidPrefix(t *testing.T) {
	s := resp.SimpleString{}

	input := "(OK\r\n"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, 1, c) // Read the operation
}

func TestSimpleString_ReadFrom_InvalidDelimiter(t *testing.T) {
	s := resp.SimpleString{}

	input := "+OK\r"
	c, err := s.ReadFrom(strings.NewReader(input))
	require.Error(t, err)
	require.EqualValues(t, len(input), c)
}

func TestSimpleString_ReadFrom(t *testing.T) {
	s := resp.SimpleString{}

	input := "+OK\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)

	require.Equal(t, "OK", s.String())
}

func TestSimpleString_ReadFrom_PayloadEqualToBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestSimpleString_ReadFrom_PayloadTwiceBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize*2 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestSimpleString_ReadFrom_PayloadHugeBuffer(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize*42 - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestSimpleString_ReadFrom_BigPayloadNotCompleteBuffers(t *testing.T) {
	s := resp.SimpleString{}

	payloadSize := bufSize + (bufSize / 2) - 3 // -3 because we have the Operation (+1) + \r\n (+2)
	payload := strings.Repeat("h", payloadSize)
	input := "+" + payload + "\r\n"

	c, err := s.ReadFrom(strings.NewReader(input))
	require.NoError(t, err)
	require.EqualValues(t, len(input), c)
	require.Equal(t, payload, s.String())
}

func TestSimpleString_WriteTo(t *testing.T) {
	input := "OK"
	expected := "+OK\r\n"
	s := resp.NewSimpleString(input)

	buf := &bytes.Buffer{}
	c, err := s.WriteTo(buf)
	require.NoError(t, err)
	require.EqualValues(t, len(expected), c)

	require.Equal(t, expected, buf.String())
}
