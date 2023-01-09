package server

import (
	"context"
	"ddia/src/logger"
	"ddia/src/resp"
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

type Server struct {
	logger   logger.Logger
	options  options
	addr     string
	listener net.Listener
	// quit signals if we want to keep listening for new incoming requests or not
	quit chan interface{}
	// wg is used for the Listen go routine, and for the goroutines processing each request
	wg       sync.WaitGroup
	handlers *Handlers
}

func New(handlers *Handlers, opts ...Option) *Server {
	options := &options{
		logger: logger.NewDiscard(),
		host:   "localhost",
		port:   6379,
	}
	for _, o := range opts {
		o.apply(options)
	}

	return &Server{
		logger:   options.logger,
		options:  *options,
		quit:     make(chan interface{}),
		handlers: handlers,
	}
}

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
				if ok := errors.As(err, &opError); ok {
					if opError.Temporary() {
						s.logger.Printf("[error][temporary] %v\n", err)
						continue
					}
				}

				s.logger.Printf("[error] %v\n", err)
				return // Non-Temporary error. Exiting with failure
			}
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.handleRequest(ctx, conn); err != nil {
				s.logger.Printf("[ERROR] handleRequest: %v\n", err)
			}
		}()
	}
}

func (s *Server) Stop() error {
	close(s.quit)
	err := s.listener.Close() // Close listener, thus new connections
	s.wg.Wait()               // Waiting for clients to finish
	s.logger.Println("Server stopped")
	return err
}

func (s *Server) Addr() string {
	return s.addr
}

// TODO: We should return error responses if something fails
func (s *Server) handleRequest(_ context.Context, conn net.Conn) error {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Printf("unable to close the connection")
		}
	}()

	s.logger.Printf("new connection from: %s", conn.RemoteAddr().String())
	for {
		cmd, err := s.readCommand(conn)
		if errors.Is(err, io.EOF) {
			s.logger.Printf("client %s closed the connection", conn.RemoteAddr().String())
			return nil
		}
		if err != nil {
			return fmt.Errorf("readCommand: %w", err)
		}

		if err := s.processCommand(conn, cmd); err != nil {
			return fmt.Errorf("processCommand: %w", err)
		}
	}
}

func (s *Server) readCommand(conn net.Conn) ([]string, error) {
	reader, operation, err := resp.PeakOperation(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to peak operation: %w", err)
	}

	s.logger.Printf("parsing operation: %q\n", string(operation))
	switch operation {
	case resp.ArrayOp:
		cmd, err := s.parseBulkString(reader)
		if err != nil {
			return nil, fmt.Errorf("parseBulkString: %w", err)
		}

		return cmd, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown opertion type: %q", operation))
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

func (s *Server) processCommand(conn net.Conn, cmd []string) error {
	if len(cmd) == 0 {
		return errors.New("invalid command: length 0")
	}

	switch verb := cmd[0]; strings.ToUpper(verb) {
	case resp.Ping:
		if err := s.handlers.Ping(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Ping: %w", err)
		}
	case resp.Get:
		if err := s.handlers.Get(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Get: %w", err)
		}
	case resp.Set:
		if err := s.handlers.Set(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Set: %w", err)
		}
	case resp.DBSize:
		if err := s.handlers.DBSize(conn, cmd); err != nil {
			return fmt.Errorf("handlers.DBSize: %w", err)
		}
	case resp.Del:
		if err := s.handlers.Del(conn, cmd); err != nil {
			return fmt.Errorf("handlers.Del: %w", err)
		}
	default:
		if err := s.handlers.UnknownCommand(conn, cmd[0]); err != nil {
			return fmt.Errorf("handlers.UnknownCommand: %w", err)
		}
	}

	return nil
}
