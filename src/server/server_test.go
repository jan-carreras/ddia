package server_test

import (
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStart(t *testing.T) {
	s := server.NewServer("localhost", 0)

	// TODO: This usage is pretty bad because we cannot start/stop the service gracefully
	// 	and we don't have any way to know if it's still alive
	require.NoError(t, s.Start())

	t.Logf("port listening to: %v", s.Addr())

	// TODO: This approach won't work if we're not adding support for reusing a connection
	// 	We would need to pass the s.Addr inside and manage the connection there (which makes total sense)
	cli := client.NewClient(s.Addr())

	rsp, err := cli.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))

	// TODO: Support multiple write operations
	//rsp, err = cli.Set("bye", "universe")
	//require.NoError(t, err)
	//require.Equal(t, `+OK\r\n`, string(rsp))

}
