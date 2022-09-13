package server_test

import (
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestStart(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 0)

	// TODO: This usage is pretty bad because we cannot start/stop the service gracefully
	// 	and we don't have any way to know if it's still alive
	require.NoError(t, s.Start())

	logger = log.New(os.Stdout, "[client] ", 0)
	// TODO: This approach won't work if we're not adding support for reusing a connection
	// 	We would need to pass the s.Addr inside and manage the connection there (which makes total sense)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))
}
