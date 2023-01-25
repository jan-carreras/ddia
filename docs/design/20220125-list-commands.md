List Commands
=============

# Purpose

## Overview

Implementation of all the list commands: `lindex`, `llen`, `lpop`, `lpush`, `lrange`, `lrem`, `lset`, `ltrim`, `rpop`, `rpush`.
 
## Terminology

* **Double linked list**: is a linked data structure that consists of a set of sequentially linked records called nodes. Each node contains three fields: two link fields (references to the previous and to the next node in the sequence of nodes) and one data field. The beginning and ending nodes' previous and next links, respectively, point to some kind of terminator, typically a sentinel node or null, to facilitate traversal of the list.
* [container/list](https://pkg.go.dev/container/list): Package list implements a doubly linked list in go.


# Background

We don't have list support on the project.


# Requirements

## Goals

* Implement the commands: `lindex`, `llen`, `lpop`, `lpush`, `lrange`, `lrem`, `lset`, `ltrim`, `rpop`, `rpush`
* Respect the complexity (big O) of the original Redis implementation.


## Non Goals

* Implement my own implementation of doubly-linked list
 

## Future goals

* Implement my own implementation, to avoid using empty interfaces for the data types. Either by using strings directly, or by using generics.
* Multi-level mutex
	* 	Lists have expensive operations (`lindex` -> `O(n)`). With the current mutex mechanism, the entire database would be locked while performing the `lindex` operation. That's unacceptable. Having a dedicated mutex on the List data structure itself would fix that problem, because we could release the main mutex after adquiring the list mutex. This is out of the scope for this implementation.

# Design chosen

## Use [container/list](https://pkg.go.dev/container/list)

In a given database, we can store an `atom` as kind (aka "type") list:

```go
// atom represents an indivisible datatype of a certain type
type atom struct {
	kind  kind
	value interface{}
}
```

Kind definition:

```go
const(
	[...]
	listKind = 4
	[...]
)
```


All the methods implemented in the `InMemory` storage must check that we're operating with the correct data-type:

```go
...
if err := m.assertType(key, listKind); err != nil {
	return err
}
...
```

All the elements of a list are `strings`. Altho `container/list` defines an element value as `any` (aka, empty interface):

```go
// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value any // <-------------------------------------
}

```

We'll assume that will always contain an `string`, because we're not going to be storing other types.

The library `container/list` is *not* thread-safe (as most of the Go stdlib) but shouldn't be a problem since we're already using a `mutex` when reading the `map` where the list is stored to, so only one thread can write at a time to the list.


* **Pros**: network interface pretty much matches `container/list` interface. Easy to implement.
* **Cons**: `container/list` are implemented using `any` instead of generics

## Test plan

Integration tests checking that the server responds correctly for each command.

 
# Resources

* [Redis lists](https://redis.io/docs/data-types/lists/)
* [Redis eBook: 1.2.2 Lists in Redis](https://redis.com/ebook/part-1-getting-started/chapter-1-getting-to-know-redis/1-2-what-redis-data-structures-look-like/1-2-2-lists-in-redis/)
* [Redis list command reference](https://redis.io/commands/?group=list)

 


