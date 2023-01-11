Handler Simplification
==============================

# Purpose

## Overview

Handlers are too verbose in code and must be simplified/rethink before we implement more commands.

## Terminology

* *Handler*: Logic in `handler.go` in charge of processing a command. The functions include validation of the input,
making the appropriate calls to the storage, return results or errors accordingly.

# Background

> The "why" and the existing context

The code seems to complex for what it is doing. Most handlers have very similar structure, thus the code can be simplified. Because this code is going to be repeated over 40 commands more, the sooner we simplify, the better:

```go
// Get the value of a key
//
//	GET key
//
// More: https://redis.io/commands/get/
func (h *Handlers) Get(conn net.Conn, cmd []string) error {
	if len(cmd) != 2 {
		err := resp.NewError("ERR wrong number of arguments for 'GET' command")
		if _, err := err.WriteTo(conn); err != nil {
			return fmt.Errorf("on invalid number of arguments: %w", err)
		}
		return nil
	}

	val, err := h.storage.Get(cmd[1])
	if errors.Is(err, ErrNotFound) {
		ok := resp.NewSimpleString("")
		if _, err := ok.WriteTo(conn); err != nil {
			h.logger.Printf("unable to write")
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("storage.Get: %w", err)
	}

	ok := resp.NewSimpleString(val)
	if _, err := ok.WriteTo(conn); err != nil {
		h.logger.Printf("unable to write")
	}
	return nil
}
```


# Requirements

## Goals

> Expected outcomes, both objective and subjective

* Simplify the signature of the method
* Define clearly how errors are handled (both to the client, or internal errors that need to be logged)
* Reduce boilerplate related to argument validation
* Reduce boilerplate on error generation & reporting to the client
* Reduce boilerplate on error handling (eg: `ErrNotFound` / `ErrValueNotInt`/ `ErrWrongKind` / ...)
* Simplify `processCommand` where we bind each command, to its handler

## Non Goals

> What intentionally are you not doing? Define and limit the scope

* Re-do all the Server implementation. Should focus first on the handler only. Further improvements can be done afterward.

## Future goals

> What things are out of the scope, but will be tackled in the future?

* Connection handling with the client will be improved
* Bulk commands will be implemented


# Design options

> Longest section on the spec, level of detail depending on the audience. Describe the engineering approach, include architecture diagram.
 
> Describe various options if possible. Define pros/cons on each one.

## Improvement 1: Simplify handler signature

Change signature from:

```go
func (h *Handlers) Get(conn net.Conn, cmd []string) error { ... }
```


```go
type client struct {
	conn net.Conn
	cmd  []string
}

func (h *Handlers) Get(c *client) error { ... }
```


* Pros: Allows extensibility by adding a struct as parameters, simplifies logic, `client` struct can have its own methods to reduce boilerplate
* Cons: 

## Improvement 2: Extract argument validation code

Validating the required is repetitive and all Handlers must have it.

```go
if len(cmd) != 2 {
	err := resp.NewError("ERR wrong number of arguments for 'GET' command")
	if _, err := err.WriteTo(conn); err != nil {
		return fmt.Errorf("on invalid number of arguments: %w", err)
	}
	return nil
}
```

This logic can be extracted in the `client` object by creating a `requiredArgs` function, that checks it for us.

```go
type client struct {
	conn net.Conn
	cmd  []string
}

func (c *client) requiredArgs(expectedArguments int) error {
	if len(c.cmd) == expectedArguments {
		err := resp.NewErrorWrongArguments(c.command())
		if _, err := err.WriteTo(c.conn); err != nil {
			return fmt.Errorf("NewErrorWrongArguments(%q): %w", c.command(), err)
		}
	}

	return nil
}

func (c *client) command() string {
	if len(c.cmd) == 0 {
		return ""
	}
	return c.cmd[0]
}

```

The code in the handler is simpler:

```go
if err := c.requiredArgs(2); err != nil {
	return err
}
```

## Improvement 3: Named arguments

Using indexes directly on `cmd` is confusing, because it's context dependent. It obscures the code.

```go
_, _ := h.storage.Get(cmd[1])
```

Instead we can use a variable to make it easier to read:


```go
key := cmd[1]
_, _ := h.storage.Get(key)
```

## Improvement 4: Rename `cmd` by `args`

`cmd` was named after a "command". It's defined as `[]string` thus each part is part of the command. Example:

```go
[]string{"SET", "hello", "world"}
```

Example of usage:

```go
func (h *Handlers) Set(_ net.Conn, cmd []string) error {
	// ...
	key, value := cmd[1], cmd[2]
	// ...
}
```


* Redis, in the C implementation, uses the `argv`+`argc` notation
* In Go `os.Args` is `[]string` the arguments passed by the CLI
* In Python it uses `sys.argv` notation

`args` might be more appropiate for a name, and closer to Go notation.


## Improvement 5: Common errors can be handled outside each specific handler

There are known responses to known errors. They don't need to be processed on each Handler, but can be handeled generically:

```go
func handleWellKnownErrors(c *client, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrNotFound) {
		emptyResponse := resp.NewSimpleString("")
		if err := c.writeResponse(emptyResponse); err != nil {
			return err
		}
		return nil
	} else if errors.Is(err, ErrWrongKind) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if err := c.writeResponse(err); err != nil {
			return err
		}
		return nil
	} else if errors.Is(err, ErrValueNotInt) {
		err := resp.NewError("ERR value is not an integer or out of range")
		if err := c.writeResponse(err); err != nil {
			return err
		}
		return nil
	}

	return err
}
```

The handlers should have the knowledge that returning those errors will result in those responses to the Client.

* Pros: Reduces boilerplate in the handlers, making error handling trivial
* Cons: Requires implicit understanding that those errors are being to be dealt with. If don't those errors to be handled, we need to write odd code in the Handler.



## Improvement 6: Simplify handler logic

Binding commands with their respective Handlers is being done in the exact same fashion for each one of them:


```go
func (s *Server) processCommand(conn net.Conn, cmd []string) error {
	if len(cmd) == 0 {
		return errors.New("invalid command: length 0")
	}

	switch verb := cmd[0]; strings.ToUpper(verb) {
	case resp.Ping:
		if err := s.handlers.Ping(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Ping: %w", err)
		}
	case resp.Set:
		if err := s.handlers.Set(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Set: %w", err)
		}
	case resp.DBSize:
		if err := s.handlers.DBSize(conn, cmd); err != nil {
			return fmt.Errorf("handlers.DBSize: %w", err)
		}
	[...]
```

The pattern repeats over an over again. Maybe using some code generation it could be automated.

Pros: Simplifies writing repetitive code
Cons: Reduces clarify, adds complexity (auto-generated code), we might encounter edge cases and exceptions in the future. 🚫


# Design chosen

> Describe which approach has been chosen, and why. Encouraged to include:  

To implement:

* Improvement 1: Simplify handler signature
* Improvement 2: Extract argument validation code
* Improvement 3: Named arguments
* Improvement 4: Rename cmd by args
* Improvement 5: Common errors can be handled outside each specific handler

We won't implement: 

*  Improvement 6: Simplify handler logic 🚫:  reason why the complexity added to the problem might be too big for the benefit it adds


## Test plan

Unit tests should keep working


 


