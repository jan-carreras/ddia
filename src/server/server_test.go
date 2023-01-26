package server_test

import (
	"testing"
	"time"
)

func TestServer_UnknownCommand(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("foo"), "ERR unknown command 'foo'"
	if rsp != want {
		t.Fatalf("missmatch: %q, want %q", rsp, want)
	}
}

func TestStart_GracefulShutdown(t *testing.T) {
	s := testServer(t)

	time.Sleep(10 * time.Millisecond)

	err := s.Stop()
	if err != nil {
		t.Fatalf("Stop faield: %v, wanted no error", err)
	}
}
