package storage

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type InMemory struct {
	records    map[string]string
	recordsMux *sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		records:    make(map[string]string),
		recordsMux: &sync.RWMutex{},
	}
}

func (m *InMemory) Set(key, value string) error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()
	m.records[key] = value

	return nil
}

func (m *InMemory) Get(key string) (string, error) {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()
	val, ok := m.records[key]
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}
