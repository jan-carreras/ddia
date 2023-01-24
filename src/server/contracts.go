package server

import "errors"

// ErrNotFound is to be returned when some methods to not find the Key. Each method that needs to implement
// it will have a comment on its signature. The implementations that do not return this error on the given
// condition, must be considered invalid.
var ErrNotFound = errors.New("not found")

// ErrValueNotInt is to be used if we're using integer operations on keys that do not hold an integer value
// Eg:
//
//		 set "hello" "world"
//	     incrby "hello" 10 <--- This operation is invalid and should return an error
var ErrValueNotInt = errors.New("value not int")

// ErrWrongKind is used when we're performing an operation on a type that does not support it
// Eg:
//
//	set "hello" "world"
//	llen "hello" <--- List Length command is invalid on a String type. Must return ErrWrongKind
var ErrWrongKind = errors.New("wrong type")

// ErrWrongNumberArguments is thrown when a command is not called with the correct number of arguments
var ErrWrongNumberArguments = errors.New("wrong number of arguments")

// ErrDBIndexOutOfRange is thrown when SELECT {idx} and idx is less than 0, or greater than the available DBs
var ErrDBIndexOutOfRange = errors.New("db index out of range")

// ErrOperationNotPermitted is returned when the user is not authenticated and cannot perform that command
var ErrOperationNotPermitted = errors.New("operation not permitted")

// Storage defines the interface that the Server needs to store things
type Storage interface {
	atomic
	stringOperations
	genericOperations
	serverOperations
}

type atomic interface {
	Lock()
	Unlock()
}

type stringOperations interface {
	// Get returns value of the given key. If the key is not found, returns ErrNotFound
	Get(key string) (string, error)
	// Set stores or overwrites the key with the given value
	Set(key, value string) error
	// IncrementBy increments the counter key by amount, returning the new value
	IncrementBy(key string, amount int) (string, error)
	// FlushDB removes all keys in the database
	FlushDB() error
	// Exists returns ErrNotFound if key does not exist, return null otherwise
	Exists(key string) error
}

type genericOperations interface {
	// Del removes a key. Returns true if existed, False otherwise.
	Del(key string) bool
	// RandomKey return a random key from all the records on the present database
	RandomKey() (string, bool)
	// Rename renames key to newkey. It returns an error when key does not exist.
	Rename(oldKey string, newKey string) error
}

type serverOperations interface {
	// Size returns the number of keys being stored
	Size() int
}
