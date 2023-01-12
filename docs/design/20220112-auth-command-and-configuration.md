Implement AUTH command & Configuration
======================================

# Purpose

## Overview

`Auth` command authenticates the client to the server. This form just authenticates against the password set with `requirepass`. In this configuration Redis will deny any command executed by the just connected clients, unless the connection gets authenticated via AUTH.
 
## Terminology

* *Authenticated client*: A client that has executed the AUTH command successfully.
* *Password protected server*: Redis server configured with `requirepass`



# Background

We don't have a authentication system of any kind, nor a system to read configuration files. Client's that connect via TCP can run any command.


# Requirements

## Goals

* Implement a Configuration system that it's able to read a [redis.conf file](https://redis.io/docs/management/config-file/) so that we can read the `requirepass` directive.
	* Use the configuration file to read the number of databases that we want to create (default=16). Now those are hardcoded.
* Implement the `AUTH` command.

## Non goals

* Implement the `include` directive

## Future goals

* The Configuration system is going to be used for many other commands, and server config (port binding, client's timeout, tcp-keepalive, tls, logfile, loglevel, ...)


# Design chosen

## Configuration System

We need to read a file that has the following syntax:

```
# Comment

key value
loglevel verbose
save 900 1
save 300 10
save 60 10000
slave-serve-stale-data yes
requirepass foobared

include /path/to/local.conf
include /path/to/other.conf
include /path/to/fragments/*.conf
```

### Encoding format

* Commends and empty lines are to be ignored. 
	* Before checking for emptiness, we'll trim trailing spaces on each line
* `keys` match `[a-z-]+`. A key can have multiple value (see `save`). A key must always have a `value` associated with it
* `value` can be any string. The content of the value depends on the specific directive.
* `include` directive is to be used to import other configuration files

### Parsing

Since keys can have duplicated, we can read the file into a datastructure like `map[string][]string` on our `Config` object.

`include` directive is going to detect if it's a single file, and try to import directly. Will return `ErrInvalidFile` if not found.

If it includes a glob/wildcard patterns [path.Glob](https://pkg.go.dev/path/filepath#Glob) can be used to expand it. Obviously, the rules of that file must be processed where the `include` directive is placed.

Warning! Should the directives in an include path OVERWRITE previously defined directives? 

### Interface

```go
var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidType = errors.New("invalid type")

func New(configPath string) (Config, error) {...}

// Get returns the first value for key found on the fil
func (c Config) Get(key string) (value string, ok bool) {


// GetD gets the value for key, returns def if key not foun
func (c Config) GetD(key, def string) string { ... }

// Integer returns the value of key as integer. If key does not exist, return def. If not integer, returns ErrInvalidType
func (c Config) Integer(key string, def int) (int, error) { ... }


// GetM returns all the values for the given key
func (c Config) GetM(key string) (values []string, ok bool) { ... }
```


## Auth Command

If directive `requirepass` is configured in the server, we cannot accept any command from an authenticated client.

We can modify the client:

```go
type client struct {
	[...]
	authenticated bool
	[...]
}
```

So, to process any command (any but `AUTH`) we must check that `requirepass` is not nil, and that `client.authenticated` is true. If that's the case, we can execute the command. Otherwise we need to fail with `-ERR operation not permitted\r\n`.

`AUTH` command is the only one that does not require authentication, and will check that the password being sent and the one configured on `requiredpass` are equal, to et `client.authenticated` to true.

**Limitations**: In the event that the password can change with the Redis service running, we'll need to iterate thru all the clients and set to false the `authenticated` flag, so that clients re-authenticate again.


## Test plan

* We have to make sure that the configuration is read correctly, even edge cases like glob files in includes
* Test that the client cannot run any command if needs to be authenticated


 


