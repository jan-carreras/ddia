# Designing Data-Intensive Applications

I love the
book [Designing Data-Intensive Applications](https://amazon.es/Designing-Data-Intensive-Applications-Reliable-Maintainable/dp/1449373321)
.
It's one of my favourite tech books of all times.

## Current situation

* I read the book
* I truly enjoy it
* I learn or refresh something new
* I shelve the book
* I forget most of what I've learnt
* Next year I repeat the cycle (3 cycles so far...)

## Problem

* I would rather not forget the concepts, or understand them in more depth
* I almost never put in practice the concepts proposed by the book
* It gives me a false sense of progression or learning

## Possible Solutions

### Option 0: Keep the cycle

Keep reading the book once a year until...?

### Option 1: Try to memorise the concepts

I've already tried and used Anki flashcards to practice spaced repetition.
It kinda works, but it's boring after a while. I get some value out of it, tho.

### Option 2: Practice using the technologies

Install the various systems (Kafka, HBase, MongoDB, Redis, ...) and play around
them to understand how they work, and try to learn more in a practical
way.

The main pain-point with that, is that without a clear objective it's a little like
playing with sand. It has not purpose, and it's easy washed away.

### Option 3: Practice _implementing_ the technologies

The book talks about Hash Indexes, SSTables and LSM-Trees, B-Trees, compaction process,
in-memory storages, replication, consensus, Write-ahead logs...

All those concepts are very appealing to me, and I love learning about them.

## Solution Chosen

Option 3 is the most appealing to me, right now. It gives me hands-on
experience coding algorithms with a real objective, exposes weak spots
(eg: network programming), encourages me to investigate further some topics,
and to read some other's peoples code to understand how they have solved
the problems.

# Challenge / Exercise

## What

I love Redis. I know almost nothing about their internals. I barely
know all the commands that it exposes, and how the cluster is configured and maintained.

But I love it anyway. What I know about its architecture, tho, is that the creator
of Redis usually spend weeks/months thinking how to implement a new feature
in Redis so that it would be as performant as possible. He was adamant to add new
things that other solutions in the market might have already, without giving it
proper though and consideration on the implications of adding this feature.

## Why?

The APIs are very clean, the network protocol is very simple, it's a key-value storage,
but I can implement persistence, replication, fault tolerance, read-replicas, etc...

## How?

Adding a list of challenges I want to solve, or features I want to implement/clone from Redis
and try to do it.





