Append only file - Sync Every Second
==============================

# Purpose

## Overview

Implement directive `appendfsync=everysec` which means that every second we'll call `Fsync` so that the OS writes
the AOF information that might be on the OS buffers, to disk.
 
## Terminology

* **AOF**: (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the Redis protocol itself.
* **fsync**: [fsync](https://linux.die.net/man/2/fsync) is used to synchronize a file's in-core state with storage device. In other words, asking the OS to actually write the information that might have in the OS buffers, to disk.
* **Flush**: Write the buffers (might be OS buffers, or Go's `bufio` data structure) to the underlying logic.


# Background

The application supports `appendfsync=always`, which calls fsync on each write operation (and it's very, very slow), and
`appendfsync=never`, that doesn't tell the OS when to sync the OS buffer to the disk, and the OS decides that.

The problem with `always` is that  is that it's very slow (350 ops/s). The problem with `never` is that we don't
quite know when the OS will sync the information to disk, and we might lose data.


# Requirements

## Goals

* Implement `everysec` that will call Fsync is there are writes not FSynced yet

## Non Goals

## Future goals

* At some point the directive `appendfsync` might be updated after the server has started using the `CONFIG` command.


# Design chosen


## Gorountine spawned on AppendOnlyFile creation

In the case that the option is `everysec`, we can use a goroutine and [time.Ticker](https://pkg.go.dev/time#Ticker) to check if there are new changes since the last time we synced to disk.

To be able to stop the ticker, we must base it on a context passed on the constructor.

* Pros: Simple to implement. Can safely stop the ticker if needed.
* Cons: With this design we cannot change the frequency of the ticker, altho right now we don't need it.





