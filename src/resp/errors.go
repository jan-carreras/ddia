package resp

import (
	"fmt"
	"io"
)

// Error defines a protocol error
type Error struct {
	string string
}

// NewError returns an Error
func NewError(string string) *Error {
	return &Error{string: string}
}

// ReadFrom reads from the Reader and loads the Error object
// Example: "+OK\r\n"
func (s *Error) ReadFrom(r io.Reader) (readCount int64, err error) {
	err = checkOperation(r, ErrorOp)
	readCount = 1 // Read the first byte
	if err != nil {
		return readCount, err
	}

	c, str, err := readFrom(r)
	readCount += c
	if err != nil {
		return readCount, fmt.Errorf("readFrom: %w", err)
	}

	s.string = str

	return readCount, nil
}

// WriteTo writes the information on Error and dumps it into the Writer
func (s *Error) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c%s\r\n", byte(ErrorOp), s.string)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

// String returns the String representation of the object
func (s *Error) String() string {
	return s.string
}

// Bytes returns the String representation encoded in []bytes
func (s *Error) Bytes() []byte {
	return []byte(s.string)
}
