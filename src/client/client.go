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
func (c *Client) response(socket net.Conn) ([]byte, error) {
	c.logger.Printf("waiting for response...")

	buf, err := io.ReadAll(socket)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf(`unable to read until delimiter \n: %w`, err)
	}
	c.logger.Printf("got response")
	return buf, nil
}

// encodeBulkStrings encodes a slice of strings into a RESP Array consisting only Bulk Strings
func encodeBulkStrings(cmd []string) ([]byte, error) {
	bulkStr := resp.NewBulkStr(cmd)
	buf := &bytes.Buffer{}
	_, err := bulkStr.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
