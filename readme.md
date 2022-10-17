# Designing Data-Intensive Applications

This repository is a playground where I implement ideas from the book
Designing Data-Intensive Applications.

## TODO

* [x] Gracefully shutdown of the server
* [ ] Do not close the connection on each command
* [ ] Be able to define read/write timeouts on the server side to prevent DoS
* [ ] Store values in-memory in a key-value store (map + mutexes)
* [ ] Benchmark the key value store and compare it with Redis
* [ ] Create my own hashmap implementation
* [ ] Be able to persist data even after graceful restart
* [ ] Be able to persist data even after a kill -9
* [ ] Be able to have a read-replica
* [ ] Implement a heartbeat system to monitor the status of the system
* [ ] Allow client to define how many replicas must commit the information before OK
