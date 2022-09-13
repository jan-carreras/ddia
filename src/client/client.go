package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Client struct {
	socket io.ReadWriteCloser
}

func NewClient(s io.ReadWriteCloser) *Client {
	return &Client{socket: s}
}

func (c *Client) Set(key, value string) ([]byte, error) {
	cmd, err := c.encodeStr([]string{"SET", key, value})
	if err != nil {
		return nil, fmt.Errorf("unable to encode the string: %w", err)
	}

	if err := c.send(cmd); err != nil {
		return nil, fmt.Errorf("unable to send: %w", err)
	}

	if rsp, err := c.recv(); err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	} else {
		return rsp, nil
	}

}

func (c *Client) send(msg []byte) error {
	if _, err := c.socket.Write(msg); err != nil {
		return fmt.Errorf("unable to write to socket: %w", err)
	}

	return nil
}

func (c *Client) recv() ([]byte, error) {
	buf, err := bufio.NewReader(c.socket).ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf(`unable to read until delimiter \n: %w`, err)
	}
	return buf, nil
}

// encodeStr encodes a slice of strings into a RESP Array consisting only Bulk Strings
func (c *Client) encodeStr(cmd []string) ([]byte, error) {
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
