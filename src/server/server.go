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
	port   int // TODO: Should this be an integer?
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

	if _, err := net.ResolveTCPAddr(serverNetwork, listen.Addr().String()); err != nil {
		return fmt.Errorf("unable to parse address: %w", err)
	}

	defer func() { _ = listen.Close() }()

	for {
		conn, err := listen.Accept()
		if err != nil {
			// TODO: What about the other connections?
			//  	What are the possible failures?
			return fmt.Errorf("listen.Accept: %w", err)
		}

		go handleRequest(conn)
	}
}

func (s *Server) TCPAddr() *net.TCPAddr {
	// TODO: Already checked before
	addr, _ := net.ResolveTCPAddr(serverNetwork, s.listen.Addr().String())
	return addr
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

	// TODO: What happens if we receive and EOF?
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
