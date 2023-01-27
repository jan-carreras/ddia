Expiring keys
=============

# Purpose

## Overview

We want to support automatic key expiration. `expire` is used to set a Time Time Live on a given key. `ttl` is used
to read what's the TTL on a given key. `persist` (later on included, not in version 1.0) is to remove the TTL on any
key.

Set a timeout on key. After the timeout has expired, the key will automatically be deleted. A key with an associated
timeout is often said to be volatile in Redis terminology.

 
## Terminology

* **TTL** or Time To Live: A fix point in time when the object becomes stale and should not be valid. In our case, a
key _after_ the TTL should be considered as "deleted".

# Background

The only way to remove a key from Redis is to use `DEL`.

# Requirements

## Goals

> Expected outcomes, both objective and subjective

* Be able to use the [expire](https://redis.io/commands/expire/) command `EXPIRE key seconds` and the key will be removed after the given time.
* Implement [TTL](https://redis.io/commands/ttl/) which will return the number of seconds remaining for the key to expire. Will return `-1` if the key
exists but has no expiration assigned. Will return `-2` if the key does not exist.
* Doing a `set` operation (or alike) on the key removes the previous TTL that it might have.
* Two operation modes: 
	* `lazy`: Read the expiration date of all the objects when being read. If the object is stale, pretend that it does not exist.
	* `active`: Actively remove stale keys and release memory from it.


## Non Goals

* Implement the `NX`, `XX`, `GT` and `LT` command arguments that were implemented in the Redis 7.0
* Clarification: Key expiration works on key objects. It does **not** work on values of complex datatypes. Eg: a specific element of a List). You can only set a TTL on the entire list.


# Design options

> The longest section on the spec, level of detail depending on the audience. Describe the engineering approach, include architecture diagram.
 
> Describe various options if possible. Define pros/cons on each one.

## Extending atom struct

Extending `atom` adding a UNIX timestamp should be enough to implement our `lazy` operation.


```go
// atom represents an indivisible datatype of a certain type
type atom struct {
	kind      kind
	value     any
	expiresAt int64 // representing a UNIX timestamp. zero value means "no expiration time"
}

func (a atom) IsExpired() bool {
	return a.expiresAt < time.Now().Unix()
}
```

## Change accessors to keys

We cannot approach "getting" a key this way, but we need a method that will check for the `expiresAt`.

```go
a, ok := m.records[key]
if !ok {
	return "", server.ErrNotFound
}
```

Needs to be access using this proxy function:

```go

func (m *InMemory) get(key string) (atom, error) {
	a, found := m.records[key]
	if !found || a.IsExpired() {
		return atom{}, server.ErrNotFound
	}

	return a, nil
}
```

## Atomate invalidation

We need to use a datastructure that can be very fast fetching elements that are the next to expire. A `heap` can do the trick in our case. Go has [container/heap](https://pkg.go.dev/container/heap) in the standard library and we can use it.

The idea would be 


* **Pros**:
* **Cons**:

# Design chosen

> Describe which approach has been chosen, and why. Encouraged to include: 

* Data Model: Schema definitions, New data models, Modified data models, Data validation methods
* Business Logic: API changes, Pseudocode, Flowcharts, Error states, Failure scenarios, Conditions that lead to errors and failures, Limitations
* Other questions: How will the solution scale? What are the limitations of the solution? How will it recover in the event of a failure? How will it cope with future requirements?


## Test plan

> How we're going to make sure that everything works as expected (unit/integration/QA)

## Monitoring and alerting

> Logging, alerting, new metrics, ...

 


# Resources

> Links, resources, useful information, etc...

 


