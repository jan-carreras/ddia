package server

import "errors"

var ErrNotFound = errors.New("not found")

type Storage interface {
	// Get returns value of the given key. If the key is not found, returns ErrNotFound
	Get(key string) (string, error)
	// Set stores or overwrites the key with the given value
	Set(key, value string) error
	// Size returns the number of keys being stored
	Size() int
	// Del removes a key, returing true if existed. False otherwise.
	Del(key string) bool
}
