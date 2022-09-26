package server

import (
	"ddia/src/logger"
	"ddia/src/resp"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	serverNetwork = "tcp"
	responseOK    = `+OK\r\n`
)

type Server struct {
	logger logger.Logger
	host   string
	port   int
	addr   string
	listen net.Listener
}

func NewServer(logger logger.Logger, host string, port int) *Server {
	return &Server{logger: logger, host: host, port: port}
}

func (s *Server) Start() error {
	listen, err := net.Listen(serverNetwork, fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	s.listen = listen
	s.addr = listen.Addr().String()

	s.logger.Printf("listening at %q...", s.addr)

	go func() {
		defer func() {
			// TODO: Implement graceful shutdown
			_ = listen.Close()
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				s.logger.Printf("[error] %v\n", err)
				continue
			}

			go func() {
				if err := s.handleRequest(conn); err != nil {
					s.logger.Printf("[ERROR] handleRequest: %v\n", err)
				}
			}()
		}
	}()

	return nil
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) handleRequest(conn net.Conn) error {
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

	return nil
}
