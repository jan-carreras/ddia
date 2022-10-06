package server_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 0, storage.NewInMemory())

	require.NoError(t, s.Start(context.Background()))

	logger = log.New(os.Stdout, "[client] ", 0)
	// TODO: This approach won't work if we're not adding support for reusing a connection
	// 	We would need to pass the s.Addr inside and manage the connection there (which makes total sense)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))
}

func TestServer_Set(t *testing.T) {
	store := storage.NewInMemory()

	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 0, store)

	require.NoError(t, s.Start(context.Background()))

	logger = log.New(os.Stdout, "[client] ", 0)
	cli := client.NewClient(logger, s.Addr())

	rsp, err := cli.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, `+OK\r\n`, string(rsp))

	v, err := store.Get("hello")
	require.NoError(t, err)
	require.Equal(t, "world", v)
}

func TestStart_GracefulShutdown(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	s := server.NewServer(logger, "localhost", 0, storage.NewInMemory())

	ctx := context.Background()

	require.NoError(t, s.Start(ctx))

	time.Sleep(100 * time.Millisecond)
	require.NoError(t, s.Stop())
}
