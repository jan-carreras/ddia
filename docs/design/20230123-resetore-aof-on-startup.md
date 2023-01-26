Restore AOF File on Startup
============================

# Purpose

## Overview

Storing a AOF file and not restoring it back it's pretty much... useless. This document describes how to import an AOF file to recover the state of the service.
 
## Terminology

## Terminology

* **AOF**: (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the Redis protocol itself.


# Background

We have mechanisms to generate an AOF file, but not a way to restore it.


# Requirements

## Goals

* Restore the AOF on startup
* Do not add new records in the AOF file while restoring the AOF
* Commands from the AOF are executed _without_ taking into account password protection in the server side


# Design options

> Longest section on the spec, level of detail depending on the audience. Describe the engineering approach, include architecture diagram.
 
> Describe various options if possible. Define pros/cons on each one.

## Option 1: Use a client to send the commands read in the AOF

Create a real Redis Client that will read the AOF file and reply the commands to the Server, validating that the server responses are OK and stopping if an error has been found.

* **Pros**: Reading logic would be independent to the server
* **Cons**: We'll need to authenticate. We'll need to tell the server that an import is being done, thus it needs to stop writing new records on the AOF file, if the Server produces an error, we need the client to stop sending information and tell the server to stop, everything needs to ocurr thru TCP, adding unnecesary overhead 🚫🚫🚫

## Option 2: On server startup, read the file and process the commands

On server startup, see if there is an AOF file available and read it, replying all the commands. In the error returns an error we'll stop processing more commands and not start the server.

Server will not listen for TCP connections until it has imported all the file.

This process will by-pass authentication by pretending that the client is authenticated.

* **Pros**: Easy to implement and fast (no network involved), prevents starting the server without a valid state
* **Cons**: server startup logic gets more complicated

# Design chosen

Option 2, since Option 1 is discarded by its own.

