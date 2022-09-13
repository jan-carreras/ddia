package client

import (
	"bytes"
	"ddia/src/logger"
	"fmt"
	"io"
	"net"
	"strconv"
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
	cmd, err := c.encodeBulkStrings([]string{"SET", key, value})
	if err != nil {
		return nil, fmt.Errorf("unable to encode the string: %w", err)
	}

	socket, err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the remote server: %w", err)
	}
	defer socket.Close()

	c.logger.Print("sending command")
	if err := c.send(socket, cmd); err != nil {
		return nil, fmt.Errorf("unable to send: %w", err)
	}

	if rsp, err := c.response(socket); err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	} else {
		return rsp, nil
	}

}

func (c *Client) connect() (net.Conn, error) {
	c.logger.Printf("connecting to %q...", c.addr)
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, fmt.Errorf("unable to dial %q: %w", c.addr, err)
	}

	return conn, nil
}

// send the msg to the socket
func (c *Client) send(socket net.Conn, msg []byte) error {
	c.logger.Print("writing on socket...")
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
func (c *Client) encodeBulkStrings(cmd []string) ([]byte, error) {
	length := len(cmd)
	if length == 0 {
		length = -1
	}

	buf := bytes.Buffer{}
	_, err := buf.WriteString(`*` + strconv.Itoa(length) + `\r\n`)
	if err != nil {
		return nil, fmt.Errorf("unable to start message: %w", err)
	}

	for _, s := range cmd {
		_, err := buf.WriteString(`$` + strconv.Itoa(len(s)) + `\r\n` + s + `\r\n`)
		if err != nil {
			return nil, fmt.Errorf("unable to write a word in message: %w", err)
		}
	}

	return buf.Bytes(), nil
}
