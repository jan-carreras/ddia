// Package aof provides persistence capabilities on the Redis server using Append Only Files.
package aof

import (
	"context"
	"io"
	"time"
)

type writeSyncer interface {
	io.WriteCloser
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
	file      writeSyncer
	options   options
	lastWrite time.Time // Only updated when option EverySecondSync is used
}

// NewAppendOnlyFile creates an AppendOnlyFile. You can pass io.Discard to the writeSyncer if you're not interested
// into saving any data.
func NewAppendOnlyFile(ctx context.Context, f writeSyncer, o options) *AppendOnlyFile {
	aof := &AppendOnlyFile{file: f, options: o}

	if o == EverySecondSync {
		go aof.startTicker(ctx)
	}

	return aof
}

func (a *AppendOnlyFile) startTicker(ctx context.Context) {
	var lastSync time.Time
	ticker := time.Tick(time.Second) // nolint: staticcheck

	for {
		select {
		case <-ticker:
			if !lastSync.Equal(a.lastWrite) {
				if err := a.file.Sync(); err != nil {
					panic(err)
				}
				lastSync = a.lastWrite
			}
		case <-ctx.Done():
			return // stop goroutine
		}
	}
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
	} else if a.options == EverySecondSync {
		a.lastWrite = time.Now()
	}

	return n, nil
}

// Close the AOF
func (a *AppendOnlyFile) Close() error {
	_ = a.file.Sync()
	return a.file.Close()
}

// NoOpAOF no-operation AOF. Everything ends up in /dev/null
type NoOpAOF struct{}

// NewNoOpAOF new NoOpAOF
func NewNoOpAOF() *NoOpAOF {
	return &NoOpAOF{}
}

// Write does nothing
func (a *NoOpAOF) Write(data []byte) (int, error) {
	return len(data), nil
}

// Close does nothing
func (a *NoOpAOF) Close() error {
	return nil
}
