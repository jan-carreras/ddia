package server_test

import (
	"ddia/src/client"
	"ddia/src/server"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	s := server.NewServer("localhost", 0)

	// TODO: This usage is pretty bad because we cannot start/stop the service gracefully
	// 	and we don't have any way to know if it's still alive
	go func() { require.NoError(t, s.Start()) }()

	time.Sleep(100 * time.Millisecond)

	t.Logf("port listening to: %v", s.TCPAddr().Port)

	conn, err := net.DialTCP("tcp", nil, s.TCPAddr())
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
