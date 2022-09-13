package server

import (
	"ddia/src/logger"
	"fmt"
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
					s.logger.Printf("[ERROR] handleRequest: %w", err)
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
		s.logger.Printf("Connection closed %s", conn.RemoteAddr().String())
	}()
	if err := conn.SetDeadline(time.Now().Add(time.Second)); err != nil {
		return fmt.Errorf("error when setting the timeout: %w", err)
	}

	// TODO: We're not processing requests with more than 1024 bytes
	buf := make([]byte, 1024) // TODO: Why 1024?

	s.logger.Printf("new connection from: %s", conn.RemoteAddr().String())

	n, err := conn.Read(buf)
	if err != nil {
		// TODO: What are the possible network errors here?! We need to know
		return fmt.Errorf("error when reading: %v", err)
	}

	if n != 0 {
		s.logger.Printf(string(buf[:n]))
	}

	// TODO: Check that we've written all the bytes we wanted to
	if _, err := conn.Write([]byte(responseOK)); err != nil {
		s.logger.Printf("unable to write")
	}

	return nil
}
