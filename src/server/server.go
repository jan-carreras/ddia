package server

import (
	"fmt"
	"log"
	"net"
)

const (
	serverNetwork = "tcp"
	responseOK    = `+OK\r\n`
)

type Server struct {
	host   string
	port   int
	addr   string
	listen net.Listener
}

func NewServer(host string, port int) *Server {
	return &Server{host: host, port: port}
}

func (s *Server) Start() error {
	listen, err := net.Listen(serverNetwork, fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	s.listen = listen
	s.addr = listen.Addr().String()

	go func() {
		defer func() {
			// TODO: Implement graceful shutdown
			_ = listen.Close()
		}()

		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("[error] %v\n", err)
				continue
			}

			go handleRequest(conn)
		}
	}()

	return nil
}

func (s *Server) Addr() string {
	return s.addr
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024) // TODO: Why 1024?

	log.Printf("New connection from: %s", conn.RemoteAddr().String())

	// TODO: We're not processing requests with more than 1024 bytes
	n, err := conn.Read(buf)
	if err != nil {
		// TODO: What are the possible network errors here?! We need to know
		log.Printf("error when reading: %v", err)
		return
	}

	// TODO: What happens if we receive an EOF?
	if n != 0 {
		log.Printf(string(buf[:n]))
	}

	// TODO: Check that we've written all the bytes we wanted to
	if _, err := conn.Write([]byte(responseOK)); err != nil {
		log.Printf("unable to write")
	}

	if err := conn.Close(); err != nil {
		log.Printf("unable to close the connection")
	}

	log.Printf("Connection closed %s", conn.RemoteAddr().String())
}
