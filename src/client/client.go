// Package client implements a Redis 1.0 client
package client

import (
	"bytes"
	"ddia/src/logger"
	"ddia/src/resp"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	dialTimeout = time.Second
)

// Client is to be used to make requests to the DDIA server
type Client struct {
	logger logger.Logger
	addr   string
}

// NewClient returns a Redis Client
func NewClient(logger logger.Logger, addr string) *Client {
	return &Client{logger: logger, addr: addr}
}

// Set sends the command SET {key} {value} to the server
func (c *Client) Set(key, value string) ([]byte, error) {
	cmd, err := encodeBulkStrings([]string{"SET", key, value})
	if err != nil {
		return nil, fmt.Errorf("unable to encode SET command: %w", err)
	}

	socket, err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the remote server: %w", err)
	}
	defer func() { _ = socket.Close() }()

	if err := c.send(socket, cmd); err != nil {
		return nil, fmt.Errorf("unable to send: %w", err)
	}

	rsp, err := c.response(socket)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}
	return rsp, nil
}

// Get returns the key being requested
func (c *Client) Get(key string) ([]byte, error) {
	cmd, err := encodeBulkStrings([]string{"GET", key})
	if err != nil {
		return nil, fmt.Errorf("unable to encode GET command: %w", err)
	}

	socket, err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connecto to remove server: %w", err)
	}
	defer func() { _ = socket.Close() }()

	if err := c.send(socket, cmd); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	rsp, err := c.response(socket)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	return rsp, nil

}

// Ping returns the key being send, or the string PONG
func (c *Client) Ping(key string) ([]byte, error) {
	cmd, err := encodeBulkStrings([]string{"PING"})
	if len(key) != 0 {
		cmd, err = encodeBulkStrings([]string{"PING", key})
	}

	if err != nil {
		return nil, fmt.Errorf("unable to encode PING command: %w", err)
	}

	socket, err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connecto to remove server: %w", err)
	}
	defer func() { _ = socket.Close() }()

	if err := c.send(socket, cmd); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	rsp, err := c.response(socket)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	return rsp, nil

}

func (c *Client) connect() (net.Conn, error) {
	c.logger.Printf("connecting to %q...", c.addr)
	conn, err := net.DialTimeout("tcp", c.addr, dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("unable to dial %q: %w", c.addr, err)
	}

	return conn, nil
}

// send the msg to the socket
func (c *Client) send(socket net.Conn, msg []byte) error {
	c.logger.Print("writing command on socket...")
	if _, err := socket.Write(msg); err != nil {
		return fmt.Errorf("unable to write to socket: %w", err)
	}
	c.logger.Print("done writing")

	return nil
}

// response reads from the TCP connection
func (c *Client) response(reader io.Reader) ([]byte, error) {
	c.logger.Printf("waiting for response...")

	reader, operation, err := resp.PeakOperation(reader)
	if err != nil {
		return nil, fmt.Errorf("PeakOperation: %w", err)
	}

	c.logger.Println("operation:", string(operation))

	switch operation {
	case resp.SimpleStringOp:
		s := resp.SimpleString{}
		_, err := s.ReadFrom(reader)
		if err != nil {
			return nil, fmt.Errorf("ReadFrom: %w", err)
		}

		return []byte(s.String()), nil
	case resp.BulkStringOp:
		bs := resp.Str{}
		_, err := bs.ReadFrom(reader)
		if err != nil {
			return nil, fmt.Errorf("ReadFrom: %w", err)
		}

		return []byte(bs.String()), nil
	}

	return nil, fmt.Errorf("unknown operation: %c", operation)
}

// encodeBulkStrings encodes a slice of strings into a RESP Array consisting only Bulk Strings
func encodeBulkStrings(cmd []string) ([]byte, error) {
	bulkStr := resp.NewArray(cmd)
	buf := &bytes.Buffer{}
	_, err := bulkStr.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
