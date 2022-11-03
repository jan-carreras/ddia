package client_test

import (
	"context"
	"ddia/src/client"
	"ddia/src/server"
	"ddia/src/storage"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestClient_Set(t *testing.T) {
	logger := log.New(os.Stdout, "[server] ", 0)
	store := storage.NewInMemory()
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	require.NoError(t, err)

	logger = log.New(os.Stdout, "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	rsp, err := c.Set("hello", "world")
	require.NoError(t, err)
	require.Equal(t, "OK", string(rsp))

	rsp, err = c.Set("chao", "universe")
	require.NoError(t, err)
	require.Equal(t, "OK", string(rsp))
}

func TestClient_Get(t *testing.T) {
	store := storage.NewInMemory()
	logger := log.New(os.Stdout, "[server] ", 0)
	handlers := server.NewHandlers(logger, store)
	s := server.NewServer(logger, "localhost", 0, handlers)

	err := s.Start(context.Background())
	require.NoError(t, err)

	logger = log.New(os.Stdout, "[client] ", 0)
	c := client.NewClient(logger, s.Addr())

	k, v := "hello", "world"
	err = store.Set(k, v)
	require.NoError(t, err)

	rsp, err := c.Get(k)
	require.NoError(t, err)
	require.Equal(t, "world", string(rsp))
}
