package server

import "errors"

// ErrNotFound is to be returned when some methods to not find the Key. Each method that needs to implement
// it will have a comment on its signature. The implementations that do not return this error on the given
// condition, must be considered invalid.
var ErrNotFound = errors.New("not found")

// Storage defines the interface that the Server needs to store things
type Storage interface {
	// Get returns value of the given key. If the key is not found, returns ErrNotFound
	Get(key string) (string, error)
	// Set stores or overwrites the key with the given value
	Set(key, value string) error
	// Size returns the number of keys being stored
	Size() int
	// Del removes a key. Returns true if existed, False otherwise.
	Del(key string) bool
}
