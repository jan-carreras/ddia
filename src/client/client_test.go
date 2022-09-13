package client_test

import (
	"ddia/src/client"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
	"time"
)

func TestClient_Set(t *testing.T) {
	server, cli := net.Pipe()
	require.NoError(t, server.SetDeadline(time.Now().Add(100*time.Millisecond)))
	require.NoError(t, cli.SetDeadline(time.Now().Add(100*time.Millisecond)))

	go func() {
		c := client.NewClient(cli)

		rsp, err := c.Set("hello", "world")
		require.NoError(t, err)
		require.Equal(t, []byte(`+OK\r\n`), rsp)

		//require.NoError(t, cli.Close())
	}()

	sent := make([]byte, 1024)
	n, err := server.Read(sent)
	require.NoError(t, err)
	require.Equal(t, `*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n`, string(sent[:n]))

	fmt.Println("[server] writing response OK")
	_, err = io.WriteString(server, `+OK\r\n`)
	require.NoError(t, err)
	require.NoError(t, server.Close())

	time.Sleep(100 * time.Millisecond)
}
