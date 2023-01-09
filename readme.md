# Designing Data-Intensive Applications

This repository is a playground where I implement ideas from the book Designing Data-Intensive Applications.

## Objective

Implement Redis 1.0 command set (see [`commands.md`]) for educative purposes.
[`commands.md`](https://github.com/jan-carreras/ddia/blob/master/commands.md)


## Design

This project will implement both a Client and a Server using the [Redis Serialization Protocol].

Backwards compatibility: Both the client and the server must be compatible with original Redis client/server.

[Redis Serialization Protocol](https://redis.io/docs/reference/protocol-spec/)

## TODO

* [x] Gracefully shutdown of the server
* [ ] Do not close the connection on each command
* [ ] Be able to define read/write timeouts on the server side to prevent DoS
* [x] Store values in-memory in a key-value store (map + mutexes)
* [ ] Benchmark the key value store and compare it with Redis
* [ ] Create my own hashmap implementation
* [ ] Be able to persist data even after graceful restart
* [ ] Be able to persist data even after a kill -9
* [ ] Be able to have a read-replica
* [ ] Implement a heartbeat system to monitor the status of the system
* [ ] Allow client to define how many replicas must commit the information before OK

## Commands

* [x] Get
* [x] Set
* [x] Ping