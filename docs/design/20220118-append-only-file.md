Append Only File
================

# Purpose

## Overview

Redis supports [two mechanisms for persistence](https://redis.io/docs/management/persistence/), RDB and AOF. RDB (Redis Database) create a point-in-time snapshots. AOF (Append Only File) persists a log with all the operations received by the server. Those operations can be replayed at server startup.

This document describes how we would implement an AOF mechanism for persistence.

 
## Terminology

* **RDB** (Redis Database): RDB persistence performs point-in-time snapshots of your dataset at specified intervals.
* **AOF** (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the Redis protocol itself.
* **Append-only Log**: file where we write always at the end. No modifications are allowed.
* **fsync**: [fsync](https://linux.die.net/man/2/fsync) is used to synchronize a file's in-core state with storage device. In other words, asking the OS to actually write the information that might have in the OS buffers, to disk.
* **Flush**: Write the buffers (might be OS buffers, or Go's `bufio` data structure) to the underlying logic.
* **Difference with `Flush` & `fsync`**: We need to `Flush` our `bufio` datastructures to make the actual `os.Write` call (otherwise everything is in-memory in go). Then, we can `fsync` by using `fd.Sync()` to make sure that the operating system write the information into the disk, rather than keeping it in their buffers.
* **Golden Files**: A golden file is the expected output of test, stored as a separate file rather than as a string literal inside the test code


# Background

Right now all the operations are being done in-memory, thus not persisted in the disk. On restart or crash all the information is lost.

The underlying in-memory implementation uses a Go map to store the key-values.


# Requirements

## Goals

* Be able to restart the server gracefully and have the previous state
* Be able to recover from a crash of the service
* Confirm to the client after changes have been saved in the disk
* Be able to configure the server to sync immediately (synchronously)

## Non Goals

> What intentionally are you not doing? Define and limit the scope

* Optimise _what_ we write into the AOF. Eg: file compaction.
* Be able to configure the AOF to be synced to disk every X seconds (`appendfsync everysec` and `appendfsync no`)
* Recreate the status of the database given some AOF. That will be done in a future Design document.

## Future goals


* File compactation, optimizations on what is stored on the AOF (Eg: instead of `INCR key` a `SET KEY {value}` can be stored, making compactation more effective.


# Design options

## Challenges/Doubts

* 1. Where do we add the Persistence Layer logic in the current project?
* 2. How do we handle concurrency? Single actor writes, or many actors can write on the AOF? How do we prevent corruption of the file if concurrency is enabled?
* 3. What's the cost of fsync in the disk on every operation?
* 4. The client has state: authentication, a database selected, etc... how do we persist those?

## Doubt 1: Where to add the Persistence Layer?

### 1.1 Wrapper to `Storage` interface

Implement a Wrapper to the `Storage` interface so that stores all the operations in an append only file. It would act like the "middleware" pattern.

* **Pros**: Conceptually easy to understand what this layer is doing.
* **Cons**: The `Storage` interface will become huge. It will be a lot of code for little value. Better to have the Persistence Layer in more generic place 🚫

### 1.2 Storing Primitives

Define a set/del storage primitives, and all the operations should end un using those primitives (including `INCR` or `DECRBY`). This way, instead of recording in the AOF:

```
INCR visit:counter # 1
INCR visit:counter # 2
INCR visit:counter # 3
```

We would record:

```
SET visit:counter 1
SET visit:counter 2
SET visit:counter 3
```

In this case we're not recording the original commands performed to the Redis service, but the effect on the Storage layer, in terms of the primitives.

Thus, the commands like `Set` will be rewritten from:

```go
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
```

To:

```go
// Set stores or overwrites the key with the given value
func (m *InMemory) Set(key, value string) error {
	m.recordsMux.Lock()
	defer m.recordsMux.Unlock()

	if err := m.assertType(key, stringKind); err != nil {
		return err
	}

	v := atom{kind: stringKind, value: value}
	err := m.set(key, v) # <---- Always using the "primitive" indirection to access the storage.
	if err != nil {
		...
	}

	return nil
}
```


* **Pros**: Although we're not recording a 1-1 mapping with the incoming commands, we're recording the sides effects to the storage. Future compactation is going to be easier, since `INCR` operations cannot be compacted.
* **Cons**: Primitives in the storage must be implemented, so the storage must be AOF-aware. Primitive usage is going to be a CONVENTION right now, rather than enforced by a new level or indirection.


## What do we store in the file?

The same as Redis is doing: store the [RESP](https://redis.io/docs/reference/protocol-spec/) representation in the disk.

## How to know at which database do we have to write?

We could add a `SELECT N` before each command, where `N` is the selected database for the client. Alternatively, the writer could store state of what's the latest database it has written to, and do not add the `SELECT N` statement if the database has not changed, because it would be an ineffectual statement. 

* **Pros**: It reduces the size of the file, specially if the clients never change the database.
* **Cons**: It's premature optimization and complicates the store mechanism.
* **Conclusion**: Let's add the `database N` before each command. We'll optimise it afterwards.

## Doubt 2: Concurrency model

A single writer should be doing the work. It will become the bottleneck and I need to quantify how much is it going to slow down the entire server. Investigation/Prototype needed!

## Doubt 3: Syncing to disk

We might be tempted to use buffers such as `bufio.Writer` that batches writes to disk to increase performance. This is to be avoided, since in the event of a crash we would lose whatever is in the buffer, which would be very inconvenient.

If we want to use something like `bufio.Writer` we must call `Flush` after each operation to make sure that the information gets written to disk.

### Detailed explanation

Example of the usage:


```go
	walFD, _ := os.Open("dump.aof")

	w := bufio.NewWriter(walFd)
	_, err := w.WriteString("database 0;")
	// error handling
	_, err = w.WriteString("set example 1")
	// error handling
	w.Flush()
	walFD.Sync()
```

This would only perform one `Write` system call with both operations, which is convenient.

See the implementation of `os.Sync`:

```go
// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.
func (f *File) Sync() error {
	if err := f.checkValid("sync"); err != nil {
		return err
	}
	if e := f.pfd.Fsync(); e != nil {
		return f.wrapErr("sync", e)
	}
	return nil
}
```

In turn that calls `Fsync` which, specifically for `darwin` architecture:

```go
// Fsync invokes SYS_FCNTL with SYS_FULLFSYNC because
// on OS X, SYS_FSYNC doesn't fully flush contents to disk.
// See Issue #26650 as well as the man page for fsync on OS X.
func (fd *FD) Fsync() error {
	if err := fd.incref(); err != nil {
		return err
	}
	defer fd.decref()
	return ignoringEINTR(func() error {
		_, err := fcntl(fd.Sysfd, syscall.F_FULLFSYNC, 0)
		return err
	})
}
```

**In any case**, whatever the specifics of the implementation being chosen, we need to `Flush` + `Sync` all the information after each operation.

* **Pros**: We minimise the risk of loosing data in case of `panic`/`kill -9`
* **Cons**: Will make everything slower

### Performance considerations

Writing to disk and doing an fsyncing afterwards can be very expensive. Let's define what "very expensive means" doing a benchmark:

Let's say that we want to store two commands on the AOF:

* `select 0`
* `set key:__rand_int__ VXK`

(Note: we're not serializing those string in RESP protocol for simplicify of the example)

I've come up with three scenarios:

* `fd.Write()` on each command without doing a `fd.Sync()`. The OS is going to Sync whenever it feels like (we might have data lost)
* `fd.Write()` on each command doing a `fd.Sync()`

Output on a MacBook Air M1:

```
goos: darwin
goarch: amd64
pkg: ddia/playground/aof
cpu: VirtualApple @ 2.50GHz

Benchmark_WithSync                           300           3955489 ns/op
Benchmark_WithoutSync                     354434              3480 ns/op
```

We have a 3 order of magnitude difference, which is astonishing. The difference is so huge, that I doubt I might be doing something wrong. Let's benchmark the performance of a Redis Server with the following options on:

```
appendonly yes
appendfsync always
```


```bash
$ redis-benchmark -c 10 -n 10000 -t set

Summary:
  throughput summary: 1259.92 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        7.832     2.864     7.903     8.175    10.975    18.079
```

Redis is doing 1259 write ops/sec, and from Go we're doing 300 ops/sec. There is a 4x discrepancy, which is a lot but it's not various order of magnitude difference, which leads me to believe that my benchmark is not fundamentally wrong, but just that `fd.Sync()` operations are just very expensive.

* **Conclusion**: The implementation of `appendfsync everysec` and `appendfsync no` make total sense, as we expect to see an increase on performance by 3 orders of magnitude with the possibility of data-loss.

## Doubt 4: Ignore read-only operations

We need to identify commands to be READ operation, or WRITE operations. The read operations should never be written on the AOF file for obvious reasons: they don't change the state of the storage, thus they are not needed for the recreation of the data.

* **Pros**: Smaller AOF
* **Cons**:

## Configuration directives

### `appendfsync`

We're going to support a sub-set of the equivalent Redis configuration described in [Redis: Persistence](https://redis.io/docs/management/persistence/).

* **appendfsync always**: fsync every time new commands are appended to the AOF. Very, very slow, very safe.
* **appendfsync everysec**: fsync every second. Fast enough, and you may lose 1 second of data if there is a disaster. 🚫
* **appendfsync no**: Never fsync, just put your data in the hands of the Operating System. The faster and less safe method. Normally Linux will flush data every 30 seconds with this configuration, but it's up to the kernel's exact tuning. 🚫

We're only going to be implementing the `always` directive in this proposal.

### `appendonly`

To turn on AOF, we'll use the directive `appendonly yes`

### `appenddirname`

The `appenddirname` configuration directive specifies the path where we're going to be writing our AOF files.


# Design chosen

Everything described in the `Design options` section, but the section `1.1 Wrapper to `Storage` interface`, that has been discarded because it's too complex with very little benefit.

## Tasks

1. Create the "primitive" abstraction to access the data layer.
2. Proof of concept of storing thousands of commands and doing a Flush + Sync each time. Benchmark the results and document it.


## Test plan


### Unit tests

To test that the AOF works as expected we'll execute a number of commands in sequence and capture the data that would have been written in the disk. We'll compare that to "golden files".


# Resources 

* [Issue: os: File.Sync() for Darwin should use the F_FULLFSYNC fcntl](https://github.com/golang/go/issues/26650)

 


