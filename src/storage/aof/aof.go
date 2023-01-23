// Package aof provides persistence capabilities on the Redis server using Append Only Files.
package aof

import (
	"io"
)

type writeSyncer interface {
	io.Writer
	Sync() error
}

type options int

const (
	// AlwaysSync tells the OS to write its buffers into the physical disk. It's a
	// very, very slow mode
	AlwaysSync options = iota
	// EverySecondSync calls fs.Sync() at maximum once every second. You can lose 1
	// second worth of data in some cases
	EverySecondSync options = iota
	// NeverSync never calls fs.Sync() and leaves the OS to manage when to write its
	// buffers to disk. Can have data loss in some cases
	NeverSync options = iota
)

// AppendOnlyFile stores the commands being executed in the Redis server into a
// file. It allows various disk synchronization mechanisms
type AppendOnlyFile struct {
	file    writeSyncer
	options options
}

// NewAppendOnlyFile creates an AppendOnlyFile. You can pass io.Discard to the writeSyncer if you're not interested
// into saving any data.
func NewAppendOnlyFile(f writeSyncer, o options) *AppendOnlyFile {
	if o == EverySecondSync || o == NeverSync {
		panic("not supported")
	}

	return &AppendOnlyFile{file: f, options: o}
}

// Write data into the AOF file
func (a *AppendOnlyFile) Write(data []byte) (int, error) {
	n, err := a.file.Write(data)
	if err != nil {
		return n, err
	}

	if a.options == AlwaysSync {
		if err := a.file.Sync(); err != nil {
			return n, err
		}
	}

	return n, nil
}
