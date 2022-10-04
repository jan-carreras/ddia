package client_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestClient_Set(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 0)

	err := s.Start(context.Background())
	require.NoError(t, err)

	logger = log.New(os.Stdout, "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	rsp, err := c.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))

	rsp, err = c.Set("chao", "universe")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))
}
