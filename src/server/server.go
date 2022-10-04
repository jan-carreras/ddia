package server

import (
	"context"
	"ddia/src/logger"
	"ddia/src/resp"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	serverNetwork = "tcp"
	responseOK    = `+OK\r\n`
)

type Server struct {
	logger   logger.Logger
	host     string
	port     int
	addr     string
	listener net.Listener
	// quit signals if we want to keep listening for new incoming requests or not
	quit chan interface{}
	// wg is used for the Listen go routine, and for the goroutines processing each request
	wg sync.WaitGroup
}

func NewServer(logger logger.Logger, host string, port int) *Server {
	return &Server{logger: logger, host: host, port: port, quit: make(chan interface{})}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen(serverNetwork, fmt.Sprintf("%s:%d", s.host, s.port))
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

func (s *Server) handleRequest(_ context.Context, conn net.Conn) error {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Printf("unable to close the connection")
			return
		}
		s.logger.Printf("connection closed %s", conn.RemoteAddr().String())
	}()
	if err := conn.SetDeadline(time.Now().Add(time.Second)); err != nil {
		return fmt.Errorf("error when setting the timeout: %w", err)
	}

	s.logger.Printf("new connection from: %s", conn.RemoteAddr().String())

	reader, operation, err := resp.PeakOperation(conn)
	if err != nil {
		return fmt.Errorf("unable to peak operation: %w", err)
	}

	s.logger.Printf("parsing operation: %q\n", string(operation))
	switch string(operation) {
	case `*`:
		if err := s.parseBulkString(reader); err != nil {
			return fmt.Errorf("parseBulkString: %w", err)
		}
	default:
		return errors.New(fmt.Sprintf("unknown opertion type: %q", operation))
	}

	if _, err := conn.Write([]byte(responseOK)); err != nil {
		s.logger.Printf("unable to write")
	}
	return nil
}

func (s *Server) parseBulkString(conn io.Reader) error {
	s.logger.Print("about to start parsing")
	b := resp.BulkStr{}

	if _, err := b.ReadFrom(conn); err != nil {
		return err
	}

	fmt.Println("Command sent", b.String())

	return nil
}
