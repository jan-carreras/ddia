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

// Storage defines the interface that the Server needs to store things
type Storage interface {
	stringOperations
	genericOperations
	serverOperations
}

type stringOperations interface {
	// Get returns value of the given key. If the key is not found, returns ErrNotFound
	Get(key string) (string, error)
	// Set stores or overwrites the key with the given value
	Set(key, value string) error
	// IncrementBy increments the counter key by amount, returning the new value
	IncrementBy(key string, amount int) (string, error)
	// Increment increments the counter key by 1, returning the new value
	Increment(key string) (string, error)
	// DecrementBy decrements the counter key by amount, returning the new value
	DecrementBy(key string, amount int) (string, error)
	// Decrement decrements the counter key by 1, returning the new value
	Decrement(key string) (string, error)
}

type genericOperations interface {
	// Del removes a key. Returns true if existed, False otherwise.
	Del(key string) bool
}

type serverOperations interface {
	// Size returns the number of keys being stored
	Size() int
}
