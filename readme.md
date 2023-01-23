# Designing Data-Intensive Applications

This repository is a playground where I implement ideas from the book Designing Data-Intensive Applications.

## Objective

Implement Redis 1.0 command set (see [`commands.md`](https://github.com/jan-carreras/ddia/blob/master/commands.md))
for educative purposes.


## Design

This project will implement ~~both a Client and~~ a Redis Server using the
[Redis Serialization Protocol](https://redis.io/docs/reference/protocol-spec/).

Backwards compatibility: The server must be able to "talk" with original `redis-cli` and `redis-benchmark` tools,
behaving as a real server.


## Objectives

* Setup
    * [x] Gracefully shutdown of the server
    * [x] Do not close the connection on each command
    * [x] Store values in-memory in a key-value store (map + mutexes)
    * [x] Benchmark the key value store and compare it with Redis
    * [ ] OOM management. Invalidate older keys or swap to disk
* Storage
    * [x] Persist using Append Only File
    * [x] Write to the WAL before sending OK confirmation to the client
    * [ ] Persist using point-in-time Snapshots
    * [ ] Be able to restart the server and keep the state (even after crash)
* Features
    * [ ] Implement `expire` commands (set a TTL for a key)
* Replication
    * [ ] Read-Only replica support
* [ ] TTL: Implement expiration mechanism

### Secondary objectives

* [ ] Create my own hashmap implementation
* [ ] Be able to define read/write timeouts on the server side to prevent DoS
* [ ] Implement a heartbeat system to monitor the status of the system

## Commands

See commands.md to know what commands have been implemented and what's the progress.
