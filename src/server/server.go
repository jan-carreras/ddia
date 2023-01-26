// Package server implements a Redis 1.0 server
package server

import (
	"context"
	"ddia/src/logger"
	"ddia/src/resp"
	"ddia/src/server/config"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

const (
	serverNetwork = "tcp"
)

// Server defines a Redis Server
type Server struct {
	logger   logger.Logger
	options  options
	addr     string
	listener net.Listener
	// quit signals if we want to keep listening for new incoming requests or not
	quit     chan interface{}
	quitOnce sync.Once
	// wg is used for the Listen go routine, and for the goroutines processing each request
	wg       sync.WaitGroup
	handlers *Handlers
	config   config.Config
}

// New returns a new Redis Server configured with the Options provided
func New(handlers *Handlers, opts ...Option) (*Server, error) {
	// Default options
	options := &options{
		logger: logger.NewDiscard(),
		host:   "localhost",
		port:   6379,
	}
	for _, o := range opts {
		o.apply(options)
	}

	c := config.NewEmpty()
	if options.configurationFile != "" {
		var err error
		c, err = config.New(options.configurationFile)
		if err != nil {
			return nil, err
		}

		options.password = c.GetD("requirepass", "")
	}

	return &Server{
		logger:   options.logger,
		options:  *options,
		quit:     make(chan interface{}),
		handlers: handlers,
		config:   c,
	}, nil
}

// Start starts the redis server
func (s *Server) Start(ctx context.Context) error {
	if err := s.restoreAOF(ctx); err != nil {
		return err
	}

	listener, err := net.Listen(serverNetwork, fmt.Sprintf("%s:%d", s.options.host, s.options.port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	s.listener = listener
	s.addr = listener.Addr().String()

	s.logger.Printf("Listening at %q", s.addr)

	s.wg.Add(1)
	go s.serve(ctx)

	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	s.quitOnce.Do(func() {
		close(s.quit)
	})
	err := s.listener.Close() // Close listener, thus new connections
	s.wg.Wait()               // Waiting for clients to finish
	return err
}

// Addr returns the address where the server is listening
//
//	Example: "192.0.2.1:25", "[2001:db8::1]:80"
func (s *Server) Addr() string {
	return s.addr
}

// serve is to be run as a Goroutine
func (s *Server) serve(ctx context.Context) {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.quit:
				s.logger.Printf("Gracefully shutting down...")
				return // Graceful shutdown
			default:
				var opError *net.OpError
				if errors.As(err, &opError) && opError.Temporary() {
					s.logger.Printf("[error][temporary] %v\n", err)
					continue
				}

				s.logger.Printf("[error] %v\n", err)
				return // Non-Temporary error. Exiting with failure
			}
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.logger.Printf("new connection from: %s", conn.RemoteAddr().String())

			// Initialize a client object using the connection and the default DB
			c := newClient(conn, s.options.dbs[0])
			if err := s.handleRequest(ctx, c); err != nil {
				s.logger.Printf("[ERROR] handleRequest: %v\n", err)
			}
			s.logger.Printf("connection closed: %s", conn.RemoteAddr().String())
		}()
	}
}

// handleRequest is an infinite loop that reads a command from the client, and
// process it. If the client closes the connection gracefully, or there is a
// major error, we exit closing the connection.
func (s *Server) handleRequest(_ context.Context, c *client) error {
	defer func() {
		if err := c.close(); err != nil {
			s.logger.Printf("unable to close server side connection")
		}
	}()

	for {
		err := c.readCommand()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("reading command: %w", err)
		}

		if err := s.processCommand(c); err != nil {
			return fmt.Errorf("processCommand: %w", err)
		}
	}
}

// processCommand maps/binds each command name with its handler. If the command
// is not mapped, UnknownCommand handler is called
func (s *Server) processCommand(c *client) (err error) {
	defer func() {
		// Processes all well known errors and returns a response to the client
		// accordingly
		err = handleWellKnownErrors(c, err)
	}()

	if err := s.isAuthenticated(c); err != nil {
		return err
	}

	switch strings.ToUpper(c.command()) {
	case "":
		return errors.New("invalid command: length 0")
	case Ping:
		return s.handlers.Ping(c)
	case Echo:
		return s.handlers.Echo(c)
	case Quit:
		return s.handlers.Quit(c)
	case Select:
		return s.handlers.Select(c, s.options.dbs)
	case Get:
		return s.handlers.Get(c)
	case MGet:
		return s.handlers.MGet(c)
	case SetNX:
		return s.handlers.SetNX(c)
	case Substr:
		return s.handlers.Substr(c)
	case Set:
		return s.handlers.Set(c)
	case DBSize:
		return s.handlers.DBSize(c)
	case Del:
		return s.handlers.Del(c)
	case Incr:
		return s.handlers.Incr(c)
	case IncrBy:
		return s.handlers.IncrBy(c)
	case Decr:
		return s.handlers.Decr(c)
	case DecrBy:
		return s.handlers.DecrBy(c)
	case Auth:
		return s.handlers.Auth(c, s.options.password)
	case FlushDB:
		return s.handlers.FlushDB(c)
	case FlushAll:
		return s.handlers.FlushAll(c, s.options.dbs)
	case Exists:
		return s.handlers.Exists(c)
	case Config:
		return s.handlers.Config(c, s.config)
	case RandomKey:
		return s.handlers.RandomKey(c)
	case Rename:
		return s.handlers.Rename(c)
	case LPush:
		return s.handlers.LPush(c)
	case RPush:
		return s.handlers.RPush(c)
	case LPop:
		return s.handlers.LPop(c)
	case RPop:
		return s.handlers.RPop(c)
	case LSet:
		return s.handlers.LSet(c)
	case LLen:
		return s.handlers.LLen(c)
	case LIndex:
		return s.handlers.LIndex(c)
	case LRem:
		return s.handlers.LRem(c)
	case LRange:
		return s.handlers.LRange(c)
	default:
		if err := s.handlers.UnknownCommand(c); err != nil {
			return fmt.Errorf("handlers.UnknownCommand: %w", err)
		}
	}

	return nil
}

func (s *Server) isAuthenticated(c *client) error {
	// Reasons why we consider the client to be authenticated:
	//		c.authenticated is true: Means it has previously been authenticated using a AUTH command
	//      s.options.password is empty: Means that the server does not require a password at all
	//      The command is AUTH: this command does not require authentication. How would we authenticate otherwise?
	if s.options.password == "" || c.authenticated || strings.ToUpper(c.command()) == Auth {
		return nil
	}

	return ErrOperationNotPermitted
}

// handleWellKnownErrors it's a simple way to map Go errors into "network errors"
// so that each Handler don't have to do this mapping every time. A handler, of
// course, can capture any of those exceptions and return something different if
// the particular error has a different meaning in that context (eg: EXISTS will
// return "0" if the given key is Not Found).
//
// Every time you declare a new error that has a well-known "wire format"
// remember to add it here.
//
// TODO: Possible improvement is to declare custom type errors and we can bake in
// those messages inside the error itself
func handleWellKnownErrors(c *client, err error) error {
	var rsp io.WriterTo

	if err == nil {
		return nil
	} else if errors.Is(err, ErrNotFound) {
		rsp = resp.NewSimpleString("")
	} else if errors.Is(err, ErrWrongKind) {
		rsp = resp.NewError("ERR value is not an integer or out of range")
	} else if errors.Is(err, ErrValueNotInt) {
		rsp = resp.NewError("ERR value is not an integer or out of range")
	} else if errors.Is(err, ErrWrongNumberArguments) {
		rsp = resp.NewError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", c.command()))
	} else if errors.Is(err, ErrOperationNotPermitted) {
		rsp = resp.NewError("NOAUTH Authentication required")
	} else if errors.Is(err, ErrIndexOurOfRange) {
		rsp = resp.NewError("ERR index out of range")
	}

	if rsp != nil {
		return c.writeResponse(rsp)
	}

	// Error is not a "well known" one, thus we're not going to respond anything to the client
	return err
}
