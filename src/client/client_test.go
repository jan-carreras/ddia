package client_test

import (
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_Set(t *testing.T) {
	s := server.NewServer("localhost", 0)
	err := s.Start()
	require.NoError(t, err)

	c := client.NewClient(s.Addr())

	rsp, err := c.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))

	rsp, err = c.Set("chao", "universe")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))
}
