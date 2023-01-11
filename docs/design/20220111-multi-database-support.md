Multiple database support
==============================

# Purpose

## Overview

This document describes how to add multiple database support to the project.
 
## Terminology

* Redis Database: Redis databases are a form of namespacing: all databases are still persisted in the same RDB / AOF file. However different databases can have keys with the same name, and commands like FLUSHDB, SWAPDB or RANDOMKEY work on specific databases.


# Background

Every Redis database instance will support 16 databases by default. While the default database is “0,” this can be changed to any number from 0-15 and can also be configured to support additional databases. To help avoid confusion, each database provides a distinct keyspace that is independent from all of the other databases and the database index number is listed at the end of the Redis URL.


In the current implementation we have the interface Storage that represents a single DB:

```go
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

```

# Requirements

## Goals

* Create 16 databases on startup
* Implement the [SELECT](https://redis.io/commands/select/) command so that the client can select the DB
* Implement the [FLUSHDB](https://redis.io/commands/flushdb/) and [FLUSHALL](https://redis.io/commands/flushall/) commands


## Non Goals

* Be able to modify the default number of databases via the CONFIG command, nor via a configuration file

## Future goals

* See Non-Goals; eventually we'll do those things


# Design options



## Option 1: Dependency Inversion `[]Storage` in the Handlers

Create 16 `Storage` instances on startup and pass them into the Handlers:

```go
type Handlers struct {
	logger  logger.Logger
	dbs     []Storage
}
```

The client must be updated so that it contains the `dbIdx`, by default `0`:

```go
type client struct {
	conn  net.Conn
	args  []string
	dbIdx int # <----- default 0
}
```

On each individual handler, we need to fetch the correct database and then perform an operation to it. To do that, we can create a helper function that receives the client, such as:

```go
func (h *Handlers) db(c *client) Storage {
	return h.dbs[c.dbIdx]
}
```

and would be used like:

```go
func (h *Handlers) Get(c *client) error {
	[...]

	key := c.args[1]
	value, err = h.db(c).Get(key)

	[...]
}
```


* Pros: Dependency Injection on the Handlers method, makes it easy to test. Separates the state of the client (`client.dbIdx`) from the storage itself (`handlers.dbs`)
* Cons: Code is more complex, and the notion of "selection the correct DB" is going to be repeated thruout the whole project, without any added benefit

## Option 2: Link a `Storage` object to the client

The client can be updated to contain a Storage object, which will be a pointer to a datastore anyway:

```go
type client struct {
	conn  net.Conn
	args  []string
	dbIdx int         # <--- We keep this to know where we're at
	db    Storage     # <--- We must assign a DB on every new connection or on SELECT command
}
```

Then the syntax for the handler is simplified, since each client already knows to which database it needs to write: to the one that has selected previously:


```go
func (h *Handlers) Get(c *client) error {
	[...]
	
	key := c.args[1]
	value, err := c.db.Get(key)
	
	[...]
}

```

* Pros: The syntax in the code is much clearer
* Cons: The client gets more complexity, it's not dumb anymore. Not just has state (`dbIdx`) but a pointer to the DB itself.

# Design chosen

Choosen option: `Option 2: Link a Storage object to the client`

It's much simple, and it's obvious that the `client` object is going to accumulate more and more state — and makes sense it this way.

We're going to add authentication logic soon (`AUTH` command), in the future we'll have clients for replication, etc...

Additionally, Redis implements it with the DB in the client.

`Option 1` is discarded and altho is more "purist" from the SOLID prespective (`D`: dependency injection), the tradeoff of reducing readibility and simplicity in implementation is too high for very little benefit.


 


