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

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type Server struct {
	logger   logger.Logger
	host     string
	port     int
	addr     string
	listener net.Listener
	// quit signals if we want to keep listening for new incoming requests or not
	quit chan interface{}
	// wg is used for the Listen go routine, and for the goroutines processing each request
	wg      sync.WaitGroup
	storage Storage
}

func NewServer(logger logger.Logger, host string, port int, storage Storage) *Server {
	return &Server{
		logger:  logger,
		host:    host,
		port:    port,
		quit:    make(chan interface{}),
		storage: storage,
	}
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

// TODO: We should return error responses if something fails
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

	cmd, err := s.readCommand(conn)
	if err != nil {
		return fmt.Errorf("readCommand: %w", err)
	}

	if err := s.processCommand(conn, cmd); err != nil {
		return fmt.Errorf("processCommand: %w", err)
	}

	return nil
}

func (s *Server) readCommand(conn net.Conn) ([]string, error) {
	reader, operation, err := resp.PeakOperation(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to peak operation: %w", err)
	}

	s.logger.Printf("parsing operation: %q\n", string(operation))
	switch string(operation) {
	case `*`:
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
	s.logger.Print("stargint to parse")
	b := resp.BulkStr{}

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

	switch verb := cmd[0]; verb {
	case "GET":
		if len(cmd) != 2 {
			return fmt.Errorf("get command must have 2 parts, having %d instead", len(cmd))
		}

		// TODO: Understand which can of response should we return
		// 	It's unclear to me which data-type should I return:
		//  	Simple Strings or Bulk Strings?
		_, err := s.storage.Get(cmd[1])
		if err != nil {
			return fmt.Errorf("storage.Get: %w", err)
		}

		if _, err := conn.Write([]byte(responseOK)); err != nil {
			s.logger.Printf("unable to write")
		}

	case "SET":
		if len(cmd) != 3 {
			return fmt.Errorf("set command must have 3 parts, having %d instead", len(cmd))
		}
		err := s.storage.Set(cmd[1], cmd[2])
		if err != nil {
			return fmt.Errorf("storage.Set: %w", err)
		}

		if _, err := conn.Write([]byte(responseOK)); err != nil {
			s.logger.Printf("unable to write")
		}

		return nil

	default:
		return fmt.Errorf("invalid command %q", verb)
	}

	return nil
}
