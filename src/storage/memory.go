// Package storage defines various implementations needed for the Server to run
// Each implementation must implement server.Storage
package storage

import (
	"ddia/src/server"
	"errors"
	"strconv"
	"sync"
)

// ErrTypeCorruption shouldNeverHappen™ and is being returned when an atom cannot
// be cast to its kind.
// Eg:
//
//	a = atom{kind: stringKind, value: 1234}
//	a.value.(string) // this will fail miserably, and ErrTypeCorruption should be returned
var ErrTypeCorruption = errors.New("type corruption")

// kind defines the various data types supported by the storage
type kind int

// nolint: unused
const (
	//undefinedKind defines the lack of value. It prevents the usage of the
	//zero-value of kind type as a valid type
	undefinedKind kind = 0
	// stringKind represents the String datatype
	stringKind kind = 1
	// setKind represents the String datatype
	setKind kind = 2
	// mapKind represents the String datatype
	mapKind kind = 3
)

// atom represents an indivisible datatype of a certain type
type atom struct {
	kind  kind
	value interface{}
}

func (a atom) String() (string, error) {
	if a.kind != stringKind {
		return "", server.ErrWrongKind
	}

	v, ok := a.value.(string)
	if !ok {
		return "", ErrTypeCorruption
	}
	return v, nil
}

func (a atom) Integer() (int, error) {
	v, err := a.String()
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, server.ErrValueNotInt
	}
	return i, nil
}

// InMemory is the simplest storage possible, storing everything in a Go map
type InMemory struct {
	records    map[string]atom
	recordsMux sync.RWMutex
}

// NewInMemory returns an in-memory storage
func NewInMemory() *InMemory {
	return &InMemory{
		records:    make(map[string]atom),
		recordsMux: sync.RWMutex{},
	}
}

// Set stores or overwrites the key with the given value
func (m *InMemory) Set(key, value string) error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()

	if err := m.assertType(key, stringKind); err != nil {
		return err
	}

	m.records[key] = atom{kind: stringKind, value: value}

	return nil
}

// Get returns value of the given key. If the key is not found, returns ErrNotFound
func (m *InMemory) Get(key string) (string, error) {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()

	a, ok := m.records[key]
	if !ok {
		return "", server.ErrNotFound
	}

	v, err := a.String()
	if err != nil {
		return "", err
	}

	return v, nil
}

// IncrementBy increments the counter key by amount, returning the new value
func (m *InMemory) IncrementBy(key string, amount int) (string, error) {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()

	a, ok := m.records[key]
	if !ok { // Key does not exist, we create one with default value to 0
		a = atom{kind: stringKind, value: "0"}
	}

	i, err := a.Integer()
	if err != nil {
		return "", err
	}

	i += amount // TODO: Check if the integer is out of bounds

	newValue := strconv.Itoa(i)
	a.value = newValue
	m.records[key] = a

	return newValue, nil
}

// Increment increments the counter key by 1, returning the new value
func (m *InMemory) Increment(key string) (string, error) {
	return m.IncrementBy(key, 1)
}

// Decrement decrements the counter key by 1, returning the new value
func (m *InMemory) Decrement(key string) (string, error) {
	return m.DecrementBy(key, 1)
}

// DecrementBy decrements the counter key by amount, returning the new value
func (m *InMemory) DecrementBy(key string, amount int) (string, error) {
	return m.IncrementBy(key, -amount)
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

// FlushDB removes all keys in the database
func (m *InMemory) FlushDB() error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()

	m.records = make(map[string]atom)

	return nil
}

// Exists returns ErrNotFound if key does not exists, return null otherwise
func (m *InMemory) Exists(key string) error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()

	_, ok := m.records[key]
	if !ok {
		return server.ErrNotFound
	}

	return nil
}

// assertType returns an error ErrWrongKind if the key exists, and it's different from kind
func (m *InMemory) assertType(key string, kind kind) error {
	if atom, ok := m.records[key]; ok && atom.kind != kind {
		return server.ErrWrongKind
	}
	return nil
}
