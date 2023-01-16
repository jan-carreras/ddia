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
}

// New returns a new Redis Server configured with the Options provided
func New(handlers *Handlers, opts ...Option) (*Server, error) {
	options := &options{
		logger: logger.NewDiscard(),
		host:   "localhost",
		port:   6379,
	}
	for _, o := range opts {
		o.apply(options)
	}

	if options.configurationFile != "" {
		fmt.Println(options.configurationFile, "----")
		c, err := config.New(options.configurationFile)
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
	}, nil
}

// Start starts the redis server
func (s *Server) Start(ctx context.Context) error {
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
			if err := s.handleRequest(ctx, conn); err != nil {
				s.logger.Printf("[ERROR] handleRequest: %v\n", err)
			}
		}()
	}
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

// TODO: We should return error responses if something fails
func (s *Server) handleRequest(_ context.Context, conn io.ReadWriteCloser) error {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Printf("unable to close server side connection")
		}
	}()

	// Initialize a client object using the connection and the default DB
	c := &client{conn: conn, db: s.options.dbs[0]}

	for {
		args, err := s.readCommand(conn)
		if errors.Is(err, io.EOF) {
			//s.logger.Printf("client %s closed the connection", conn.RemoteAddr().String())
			return nil
		} else if err != nil {
			return fmt.Errorf("readCommand: %w", err)
		}

		// Load the arguments to the client, to be able to process the request
		c.args = args

		if err := s.processCommand(c); err != nil {
			return fmt.Errorf("processCommand: %w", err)
		}
	}
}

func (s *Server) readCommand(conn io.Reader) ([]string, error) {
	reader, operation, err := resp.PeakOperation(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to peak operation: %w", err)
	}

	s.logger.Printf("parsing operation: %q\n", string(operation))
	switch operation {
	case resp.ArrayOp:
		args, err := s.parseBulkString(reader)
		if err != nil {
			return nil, fmt.Errorf("parseBulkString: %w", err)
		}

		return args, nil
	default:
		return nil, fmt.Errorf("unknown opertion type: %q", operation)
	}
}

func (s *Server) parseBulkString(conn io.Reader) ([]string, error) {
	s.logger.Print("starting to parse")
	b := resp.Array{}

	if _, err := b.ReadFrom(conn); err != nil {
		return nil, err
	}

	s.logger.Printf("command sent: %s\n", b.Strings())

	return b.Strings(), nil
}

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
	case resp.Ping:
		return s.handlers.Ping(c)
	case resp.Echo:
		return s.handlers.Echo(c)
	case resp.Quit:
		return s.handlers.Quit(c)
	case resp.Select:
		return s.handlers.Select(c, s.options.dbs)
	case resp.Get:
		return s.handlers.Get(c)
	case resp.Set:
		return s.handlers.Set(c)
	case resp.DBSize:
		return s.handlers.DBSize(c)
	case resp.Del:
		return s.handlers.Del(c)
	case resp.Incr:
		return s.handlers.Incr(c)
	case resp.IncrBy:
		return s.handlers.IncrBy(c)
	case resp.Decr:
		return s.handlers.Decr(c)
	case resp.DecrBy:
		return s.handlers.DecrBy(c)
	case resp.Auth:
		return s.handlers.Auth(c, s.options.password)
	case resp.FlushDB:
		return s.handlers.FlushDB(c)
	case resp.FlushAll:
		return s.handlers.FlushAll(c, s.options.dbs)
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
	if s.options.password == "" || c.authenticated || strings.ToUpper(c.command()) == resp.Auth {
		return nil
	}

	return ErrOperationNotPermitted
}

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
	}

	if rsp != nil {
		return c.writeResponse(rsp)
	}

	// Error is not a "well known" one, thus we're not going to respond anything to the client
	return err
}
