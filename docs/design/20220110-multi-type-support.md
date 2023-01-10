MultiType support on Storage
============================

# Purpose

Add support to multiple datatypes on the Storage layer.

# Background

The storage must comply with the interface defined in the server:

```go
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
```

Right know all operation are [Strings](https://redis.io/docs/data-types/strings/) data types, thus the in-memory implementation has a map and a mutex as storage:

```go
// InMemory is the simplest storage possible, storing everything in a Go map
type InMemory struct {
	records    map[string]string
	recordsMux sync.RWMutex
}
```


# Requirements

> Goals, non-goals and future goals 

## Goals 
* The goal is to be able to support datatypes like [Lists](https://redis.io/docs/data-types/lists/), [Sets](https://redis.io/docs/data-types/sets/), [Hashes](https://redis.io/docs/data-types/hashes/).
* Do it in a way that it's scalable and easy to add new data types
* Do it in an efficient way


Most Redis operations are type-specific, which means that can not be used on other datatypes. The storage must be type aware.

## Non Goals

* Implement all those datatypes

## Future goals

* In the future we'll need to persist this Storage object into disk
* In the future we'll need to replicate this Storage flow of data into the network
* Implement [EXPIRE](https://redis.io/commands/expire/) and [TTL](https://redis.io/commands/ttl/) commands.


# Design options

## Option 1: Serialized strings

One possible option is to keep the in-memory storage as is:

```go
// InMemory is the simplest storage possible, storing everything in a Go map
type InMemory struct {
	records    map[string]string
	recordsMux sync.RWMutex
}
```

and store serialized representations of each datatype: Lists, Sets, Hashes, etc...

### Pros

* Easy to implement

### Cons

* Not efficient: For hashmaps we need to deserialize the entire object to find a specific key. 🚫
* Inserts for most datatypes (eg: SortedSets) would mean to deserialize and serialize again. 🚫


## Option 2: Types

Instad of storing strings, we can store a generic datatype that includes the type. Since `type` is a reserved word in Golang, we should use another word instead.

### Atom and Kind

This object could be called `atom` (`atom` as in: indivisible part). Each `atom` would be of a given `kind` (`kind` as in `data type`, or `data kind`. We want to avoid the keyword `type`):

```go
// kind defines the various data types supported by the storage
type kind int

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
	kind kind
	// [...] it's data is to be defined
}


// InMemory is the simplest storage possible, storing everything in a Go map
type InMemory struct {
	records    map[string]atom // <--- Using the atom struct
	recordsMux sync.RWMutex
}
```
#### Sanity check

Each data type must ensure that it's performing operations to a Key with the same data type. Example of real Redis:

```redis-cli
127.0.0.1:6379> set key value
OK
127.0.0.1:6379> get key
"value"
127.0.0.1:6379> hget key f
(error) WRONGTYPE Operation against a key holding the wrong kind of value
```

To prevent accessing to the wrong key, a function like:

```go
var ErrWrongKind = errors.New("wrong type") // To be defined in the Server package

// assertType returns an error ErrWrongKind if the key exists and it's different than kind
func (m *InMemory) assertType(key string, kind kind) error {
	if atom, ok := m.records[key]; ok && atom.kind != kind {
		return ErrWrongKind
	}
	return nil
}

```

Each function on the storage that operates on a certain datatype would to a type check:

```go
func (m *InMemory) Get(key string) (string, error) {
	m.recordsMux.RLock()
	defer m.recordsMux.RUnlock()
	
	if err := m.assertType(key, stringKind); err != nil {
		return "", err
	}
	[...]
```

#### Values

If an `atom` is of certain type, we need a way to store different value **types** for each atom. How to "cast" the value into its correct type might be challenging.


##### Option 2.1: Empty interface

```go
type atomEmptyInterface struct {
	kind kind
	value interface{}
}
```

We would need to "cast" to the correct type on each usage:

```go
atom, ok := m.records[key]
// assert kind (eg: stringKind)
v, ok := atom.value.(string)
if !ok {
   // The database is corrupt. That's a pretty bad error because the atom.kind is different from atom.value type. Should never happen, but we need to check otherwise we'll panic.
}
```

* Advantages: Simple, easy to adapt and extend, minimal use of pointers
* Disadvantages: Type casting is to be avoided in Go, since the compiler cannot help. Might be source of errors.
* Conclusion: Lets use it for its simplicity ✅

##### Option 2.2: Empty interface + Storage receives Atom objects

Redis implements it this way: the storage receives the equivalent of `atom` objects and have no oppinion on them. When parsing the client's command it generates automatically a `atom` and always passes a pointer to it thru the entire stack.

It provides no abstraction between layers, but it might be the most efficient way to implement it.

* Advantages: Seems to be the most performant and ensures minimum memory copies
* Disadvantages: Breakes the "storage abstraction", and the main part of the application needs to deal with empty interfaces as a datatype. I would try to avoid it for now.
* Conlusion: I think I'm going to be end up migrating to this as a way to improve performance, but I'll delay it for now

##### Option 2.3: Embed all possible types in the `atom`

```go
type list struct{}
type set struct{}
type sortedSet struct{}

type value struct {
	str     *string
	hashmap map[string]string
	list    *list
	set     *set
	sset    *sortedSet
}

// atom represents an indivisible datatype of a certain type
type atom struct {
   kind  kind
   value value
}
```

* Advantages: Type system will work
* Disadvantages: not scalable, a lot of wasted memory for unused pointers for each key, might still have nil pointer errors
* Conclusion:  🚫 discarted



# Design choosen

The best option is to implement `Option 2: Types` as described in the previous section, with using Empty Interfaces (`Option 2.1: Empty interface`).



