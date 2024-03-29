package aof_test

import (
	"context"
	"ddia/src/storage/aof"
	"os"
	"path"
	"testing"
)

func TestWrite(t *testing.T) {
	tmpFile := path.Join(t.TempDir(), "test.aof")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	a := aof.NewAppendOnlyFile(context.Background(), f, aof.AlwaysSync)
	n, err := a.Write([]byte("some data"))
	if err != nil {
		t.Fatalf("expecting no error: %v", err)
	}

	if want := 9; n != want {
		t.Fatalf("invalid number of bytes written: %d, want %d", n, want)
	}

	n, err = a.Write([]byte("more data"))
	if err != nil {
		t.Fatalf("expecting no error: %v", err)
	}
	if want := 9; n != want {
		t.Fatalf("invalid number of bytes written: %d, want %d", n, want)
	}

	if err := a.Close(); err != nil {
		t.Fatalf("expecting no error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("expecting no error: %v", err)
	}

	if want := "some datamore data"; string(data) != want {
		t.Fatalf("information written is not correct: %q. want %q", data, want)
	}

}

func TestNoOpAOF(t *testing.T) {
	noOp := aof.NewNoOpAOF()
	n, err := noOp.Write([]byte("1234"))
	if err != nil {
		t.Fatalf("expecting no error: %v", err)
	}

	if want := 4; n != want {
		t.Fatalf("invalid count: %d expecting %d", n, want)
	}

	if err := noOp.Close(); err != nil {
		t.Fatalf("expecting no error: %v", err)
	}
}
