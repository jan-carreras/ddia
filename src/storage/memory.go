// Package storage defines various implementations needed for the Server to run
// Each implementation must implement server.Storage
package storage

import (
	"ddia/src/server"
	"sync"
)

// InMemory is the simplest storage possible, storing everything in a Go map
type InMemory struct {
	records    map[string]string
	recordsMux sync.RWMutex
}

// NewInMemory returns an in-memory storage
func NewInMemory() *InMemory {
	return &InMemory{
		records:    make(map[string]string),
		recordsMux: sync.RWMutex{},
	}
}

// Set stores or overwrites the key with the given value
func (m *InMemory) Set(key, value string) error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()
	m.records[key] = value

	return nil
}

// Get returns value of the given key. If the key is not found, returns ErrNotFound
func (m *InMemory) Get(key string) (string, error) {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()
	val, ok := m.records[key]
	if !ok {
		return "", server.ErrNotFound
	}

	return val, nil
}

// Size returns the number of keys being stored
func (m *InMemory) Size() int {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()
	return len(m.records)
}

// Del removes a key. Returns true if existed, False otherwise.
func (m *InMemory) Del(key string) bool {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()
	_, found := m.records[key]
	delete(m.records, key)
	return found
}
