# Designing Data-Intensive Applications

This repository is a playground where I implement ideas from the book Designing Data-Intensive Applications.

**NOTE 2023-01-31**: I've finished playing around with the project. In [`wrapping-up.md`](https://github.com/jan-carreras/ddia/blob/master/wrapping-up.md)
you will find the conclusions and what I've learnt. It has been a very nice project. I strongly recommend it
to anyone 🚀

## Objective

Implement Redis 1.0 command set (see [`commands.md`](https://github.com/jan-carreras/ddia/blob/master/commands.md))
for educative purposes, without any external dependency (empty go.mod).

## Design

This project will implement a Redis Server using the
[Redis Serialization Protocol](https://redis.io/docs/reference/protocol-spec/).

## Objectives

* Setup
    * [x] The server must be able to "talk" with original `redis-cli` and `redis-benchmark` tools, behaving as a real server.
    * [x] Gracefully shutdown of the server
    * [x] Do not close the connection on each command
    * [x] Store values in-memory in a key-value store (map + mutexes)
    * [x] Benchmark the key value store and compare it with Redis
    * [ ] OOM management. Invalidate older keys or swap to disk
* Storage
    * [x] Persist using Append Only File
    * [x] Write to the WAL before sending OK confirmation to the client
    * [x] Be able to restart the server and keep the state (even after crash)
    * [ ] Persist using point-in-time Snapshots
* Features
    * [x] Implement `expire` commands (set a TTL for a key)
* Replication
    * [ ] Read-Only replica support
* [x] TTL: Implement expiration mechanism

### Secondary objectives

## Commands supported

See [`commands.md`](https://github.com/jan-carreras/ddia/blob/master/commands.md) to know what commands have been implemented and what's the progress.
