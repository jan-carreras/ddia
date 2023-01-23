package server_test

import (
	"bytes"
	"context"
	"ddia/src/server"
	"ddia/src/storage/aof"
	"ddia/testing/log"
	"os"
	"path"
	"strings"
	"testing"
)

func TestServer_AppendOnlyFile(t *testing.T) {
	tmpFile := path.Join(t.TempDir(), "test.aof")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("error not expected: %v", err)
	}
	appendOnlyFile := aof.NewAppendOnlyFile(context.Background(), f, aof.AlwaysSync)

	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger, appendOnlyFile)

	s, err := server.New(handlers, serverOptions()...)
	if err != nil {
		t.Fatalf("expecting server to be able to start without problems: %v", err)
	}

	err = s.Start(context.Background())
	if err != nil {
		t.Fatalf("expecting no error: %q", err.Error())
	}

	t.Cleanup(func() { _ = s.Stop() })

	conn := testConn(t, s)

	req := func(args string) string {
		return req(t, conn, strings.Split(args, " "))
	}

	req("set key value")
	req("get key")
	req("get another-key")
	req("set second key")
	req("incrby visits 1")

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("error not expected: %v", err)
	}

	goldenFilePath := "testdata/test.aof.txt"
	/**
	// Update the golden file:

	err := os.WriteFile(goldenFilePath, content, 0600)
	if err != nil {
		t.Fatalf("unable to update golden file")
	}
	*/

	goldenFile, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatalf("expecting to be able to read the golden file: %v", err)
	}

	if !bytes.Equal(content, goldenFile) {
		t.Fatalf("golden file does not match:\n------\n%s\n------\nwant:\n------\n%s\n------", content, goldenFile)
	}

}
