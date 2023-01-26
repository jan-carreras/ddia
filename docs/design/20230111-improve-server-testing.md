Improve Server side testing
==============================

# Purpose

## Overview

Test coverage on the server implementation is deficient (less than 30%). This highlights that we either don't have a good testing strategy, or that we're just lazy. This document tries to improve the first point.

Lack of proper testing architecture is going to prevent a healthy development of the project.
 

# Background

The current strategy is deficient because it tries to do too much at the same time. The implementation of a Redis Server and a Redis Client should be independent.

Now we're doing integration tests between the Server and the Client, to make sure the both work. This is a mistake, in terms of design.

The problem being that every time that the server implements a new funcionality, the client must support it as well in order to be able to test either of them.

That's why the testing of the server is lacking behind, because the focus is to implement Server functionality. Since Client and Server are tested together, but not advancing at the same pace, we're moving at the speed of the slower development.


# Requirements

## Goals

* Be able to test the Server and Client independently of each other
* Increase the Server coverage to more than 80%
* Automatically run the tests on GitHub via a pipeline (+ format and linter)
* Do not depend on external dependencies (eg: real `redis-cli` or a go library such as `go-redis/redis`). We want to keep the project dependency-free.


## Non Goals

* We don't aim to improve the client for now, since development is going to pause to focus on the server

## Future goals

* Improve the client's testing as well, and make it independent of the server


# Design options

There are various things to be tested in the server:

1. **Commands**: Check that the commands are properly implemeneted. That's the most important part:
	1. It helps on development speed to have a short feedback loop
	1. Prevents regressions
1. **Network**: Check that the network works as expected
	1. 	Connections are closed as expected
	1. We can do concurrent requests without corrupting data
1. **Wire protocol**:
	1. Checking that the network protocol is working as expected


We need to be able to start a Server but send commands to it without going thru a socket. 


## Option 1: Commands

To test the commands starting the server as usually, but with a mechanism that bypasses the network protocol (reading sockets, parsing input, etc...).

The easiest way to test it is to start a Redis server and make requests via TCP. We'll use `resp` package to encode the messages as expected, and create some help test functions to make tests look as simple as possible.

A test can look like:

```go
func TestServer_Ping(t *testing.T) {
	req := makeReq(t)

	rsp, want := req("ping"), "+PONG\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}

	rsp, want = req("ping hello world"), "+hello world\r\n"
	if rsp != want {
		t.Fatalf("invalid response: %q want %q", rsp, want)
	}
}
```


* Pros: The code on the tests is super simple and easy to understand
* Cons: We're not de-serializing the responses, and we're checking raw RESP answers. Might be an improvement for the future.

## Option 2: Network

Connection management, graceful shutdown, correct context cancellation propagation, make sure that we're not closing the connection to a client half-way writing on it, etc...

* Pros: improves realability and stability of the solution
* Cons: too complex now to tackle now 🚫


## Option 3: Wire Protocol

84% of coverage already, which is pretty good. Won't do.

# Design chosen

* Will do: `Option 1: Commands`: Needed to keep impelementing commands reliably and a short feedback loop. It's a must to have good testing on commands.
* Won't do: `Option 2`: Testing the network is a must, but I would rather implement all the pending commands of Redis 1.0.0 and _then_ focus on network, performance, stability.
* Won't do: `Option 3`: The Wire Protocl (`resp` package) has pretty good coverage already. No need to improve it.


## Test plan

Helper functions that can help the tests:

```go

func makeReq(t *testing.T) func(string) string {
	s := testServer(t)
	conn := testConn(t, s)

	return func(args string) string {
		return req(t, conn, strings.Split(args, " "))
	}
}

func req(t *testing.T, conn net.Conn, req []string) string {
	t.Helper()

	r := resp.NewArray(req)
	_, err := r.WriteTo(conn)
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}

	buf := make([]byte, 1024) // This is going to byte my ass, for sure

	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("not expecing error: %q", err.Error())
	}

	return string(buf[:n])
}

func testConn(t *testing.T, s *server.Server) net.Conn {
	conn, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}
	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

func testServer(t *testing.T) *server.Server {
	t.Helper()

	logger := log.ServerLogger()
	handlers := server.NewHandlers(logger)

	s := server.New(handlers, serverOptions()...)
	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("expecing no error: %q", err.Error())
	}

	t.Cleanup(func() { _ = s.Stop() })
	return s
}

func serverOptions() []server.Option {
	logger := log.ServerLogger()
	dbs := make([]server.Storage, 16)
	for i := 0; i < len(dbs); i++ {
		dbs[i] = storage.NewInMemory()
	}

	return []server.Option{
		server.WithLogger(logger),
		server.WithRandomPort(),
		server.WithDBs(dbs),
	}
}
```

See on `Option 1` what what would be the usage of those helpers functions. The aim is to test all commands.
