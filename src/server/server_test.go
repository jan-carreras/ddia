package server_test

import (
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestStart(t *testing.T) {
	s := server.NewServer("localhost", 0)

	// TODO: This usage is pretty bad because we cannot start/stop the service gracefully
	// 	and we don't have any way to know if it's still alive
	require.NoError(t, s.Start())

	t.Logf("port listening to: %v", s.Addr())

	conn, err := net.Dial("tcp", s.Addr())
	require.NoError(t, err)

	cli := client.NewClient(conn)
	rsp, err := cli.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))

	// TODO: Support multiple write operations
	//rsp, err = cli.Set("bye", "universe")
	//require.NoError(t, err)
	//require.Equal(t, `+OK\r\n`, string(rsp))

}
