package resp

import (
	"fmt"
	"io"
)

var _ dataType = (*SimpleString)(nil)

// SimpleString returns a SimpleString datatype
type SimpleString struct {
	string string
}

// NewSimpleString returns a SingleString datatype
func NewSimpleString(string string) *SimpleString {
	return &SimpleString{string: string}
}

// ReadFrom reads from the Reader and loads the SimpleString object
// Example: "+OK\r\n"
func (s *SimpleString) ReadFrom(r io.Reader) (readCount int64, err error) {
	err = checkOperation(r, SimpleStringOp)
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
