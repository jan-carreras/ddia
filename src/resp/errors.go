package resp

import (
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
// Example: "-ERR unknown command\r\n"
func (s *Error) ReadFrom(r io.Reader) (readCount int64, err error) {
	readCount, s.string, err = readFrom(r)
	return readCount, err
}

// WriteTo writes the information on Error and dumps it into the Writer
func (s *Error) WriteTo(w io.Writer) (int64, error) {
	return fprintf(w, "%c%s\r\n", byte(ErrorOp), s.string)
}

// String returns the String representation of the object
func (s *Error) String() string { return s.string }

// Bytes returns the String representation encoded in []bytes
func (s *Error) Bytes() []byte { return []byte(s.string) }
