package aof

import (
	"bufio"
	"context"
	"ddia/src/resp"
	"fmt"
	"os"
)

// ImportAppendOnlyFile to be used for reading an AOF file and pretending to be a
// TCP connection so that the server can read all the commands one after the
// other. If one command fails to execute in the server, the ImportAppendOnlyFile
// returns an error afterwards.
type ImportAppendOnlyFile struct {
	ctx context.Context
	f   *os.File
	b   *bufio.Reader
	err error
}

// NewImportAppendOnlyFile returns a ImportAppendOnlyFile
func NewImportAppendOnlyFile(ctx context.Context, aofPath string) (*ImportAppendOnlyFile, error) {
	f, err := os.OpenFile(aofPath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	b := bufio.NewReaderSize(f, 2<<20) // 1MB

	return &ImportAppendOnlyFile{ctx: ctx, f: f, b: b}, nil
}

// Read is used by the Server to read the operations. It would be the equivalent
// of reading from the TCP socket.
//
// If a Server response has returned an error code, we return an error here as well.
func (i *ImportAppendOnlyFile) Read(p []byte) (n int, err error) {
	if i.err != nil {

		return 0, i.err
	}
	return i.b.Read(p)
}

// Write it's a wierd name for what we're doing here. From the point of view of
// this method we're allowing the server to write its response back to us. From a
// Redis client perspective this method should be called Read. It's important
// that we validate if the Server is reporting an error after running each
// command. If an error is found, we return an error on both the Write and Read
// calls from then on
func (i *ImportAppendOnlyFile) Write(p []byte) (n int, err error) {
	if len(p) != 0 && p[0] == resp.ErrorOp { // Error detected
		i.err = fmt.Errorf("stopping import the AOF file: %s", p[1:])
	}

	return len(p), i.err
}

// Close closes the underlying file
func (i *ImportAppendOnlyFile) Close() error {
	return i.f.Close()
}
