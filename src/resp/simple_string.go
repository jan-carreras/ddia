package resp

import (
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
func (s *SimpleString) ReadFrom(r io.Reader) (c int64, err error) {
	c, s.string, err = readFrom(r)
	return c, err
}

// WriteTo writes the information on SimpleString and dumps it into the Writer
func (s *SimpleString) WriteTo(w io.Writer) (int64, error) {
	return fprintf(w, "%c%s\r\n", byte(SimpleStringOp), s.string)
}

// String returns the String representation of the object
func (s *SimpleString) String() string { return s.string }

// Bytes returns the String representation encoded in []bytes
func (s *SimpleString) Bytes() []byte { return []byte(s.string) }
