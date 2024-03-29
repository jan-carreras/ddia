// Package storage defines various implementations needed for the Server to run
// Each implementation must implement server.Storage
package storage

import (
	"container/list"
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
	// setKind represents the List datatype
	listKind kind = 2
)

// atom represents an indivisible datatype of a certain type
type atom struct {
	kind  kind
	value any
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

func (a atom) List() (*list.List, error) {
	v, ok := a.value.(*list.List)
	if !ok {
		return nil, ErrTypeCorruption
	}
	return v, nil
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

// Lock acquires the lock in the DB. You must call Unlock after its usage
func (m *InMemory) Lock() {
	m.recordsMux.Lock()
}

// Unlock releases the lock in the DB. Calling it if the DB is not locked causes a panic
func (m *InMemory) Unlock() {
	m.recordsMux.Unlock()
}

// Set stores or overwrites the key with the given value
func (m *InMemory) Set(key, value string) error {
	if err := m.assertType(key, stringKind); err != nil {
		return err
	}

	m.records[key] = atom{kind: stringKind, value: value}

	return nil
}

// Get returns value of the given key. If the key is not found, returns ErrNotFound
func (m *InMemory) Get(key string) (string, error) {
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

// Size returns the number of keys being stored
func (m *InMemory) Size() int {
	return len(m.records)
}

// Del removes a key. Returns true if existed, False otherwise.
func (m *InMemory) Del(key string) bool {
	_, found := m.records[key]
	delete(m.records, key)
	return found
}

// FlushDB removes all keys in the database
func (m *InMemory) FlushDB() error {
	m.records = make(map[string]atom)

	return nil
}

// Exists returns ErrNotFound if key does not exist, return null otherwise
func (m *InMemory) Exists(key string) error {
	_, ok := m.records[key]
	if !ok {
		return server.ErrNotFound
	}

	return nil
}

// RandomKey return a random key from all the records on the present database
func (m *InMemory) RandomKey() (string, bool) {
	for k := range m.records {
		return k, true
	}
	return "", false
}

// Rename renames key to newkey. It returns an error when key does not exist.
func (m *InMemory) Rename(oldKey string, newKey string) error {
	value, ok := m.records[oldKey]
	if !ok {
		return server.ErrNotFound
	}
	m.records[newKey] = value
	delete(m.records, oldKey)
	return nil
}

// assertType returns an error ErrWrongKind if the key exists, and it's different from kind
func (m *InMemory) assertType(key string, kind kind) error {
	if atom, ok := m.records[key]; ok && atom.kind != kind {
		return server.ErrWrongKind
	}
	return nil
}
