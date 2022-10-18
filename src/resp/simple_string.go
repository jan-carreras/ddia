package resp

import (
	"errors"
	"fmt"
	"io"
)

var _ dataType = (*SimpleString)(nil)

type SimpleString struct {
	string string
}

func NewSimpleString(string string) *SimpleString {
	return &SimpleString{string: string}
}

// ReadFrom reads from the Reader and loads the SimpleString object
// Example: "+OK\r\n"
func (s *SimpleString) ReadFrom(r io.Reader) (readCount int64, err error) {
	err = checkOperation(r, SimpleStringOp)
	readCount += 1 // Read the first byte
	if err != nil {
		return readCount, err
	}

	c, err := s.readFrom(r)
	readCount += c
	if err != nil {
		return readCount, fmt.Errorf("readFrom: %w", err)
	}

	return readCount, nil
}

// WriteTo writes the information on SimpleString and dumps it into the Writer
func (s *SimpleString) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c%s\r\n", byte(SimpleStringOp), s.string)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

// String returns the String representation of the object
func (s *SimpleString) String() string {
	return s.string
}

// Bytes returns the String representation encoded in []bytes
func (s *SimpleString) Bytes() []byte {
	return []byte(s.string)
}

// ignoreDelimiterCharacters ignores the last two characters if they are \r\n or fails
func ignoreDelimiterCharacters(s string) (string, error) {
	if l := len(s); l < 2 {
		return "", fmt.Errorf("invalid string lenght")
	} else if s[l-2] != '\r' || s[l-1] != '\n' {
		return "", fmt.Errorf("unexpcted end")
	} else {
		s = s[:l-2] // Ignore the last two characters
	}

	return s, nil
}

func (s *SimpleString) readFrom(r io.Reader) (readCount int64, err error) {
	buf := make([]byte, readBufferSize)
	for {
		c, err := r.Read(buf)
		readCount += int64(c)

		if errors.Is(err, io.EOF) || c == 0 {
			break
		}

		if err != nil {
			return readCount, fmt.Errorf("unable to read: %w", err)
		}

		s.string += string(buf[:c])
	}

	s.string, err = ignoreDelimiterCharacters(s.string)
	if err != nil {
		return readCount, fmt.Errorf("ignoreDelimiterCharacters: %w", err)
	}

	return readCount, nil
}
