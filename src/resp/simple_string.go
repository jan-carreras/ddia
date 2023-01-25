package resp

import (
	"io"
)

var _ dataType = (*SimpleString)(nil)

// SimpleString returns a SimpleString datatype
type SimpleString struct {
	string string
	isNull bool
}

// NewSimpleString returns a SingleString datatype
func NewSimpleString(string string) *SimpleString {
	return &SimpleString{string: string}
}

// NewNullSimpleString Returns a null Simple String, represented in the network as "$-1\r\n"
func NewNullSimpleString() *SimpleString {
	return &SimpleString{isNull: true}
}

// ReadFrom reads from the Reader and loads the SimpleString object
// Example: "+OK\r\n"
func (s *SimpleString) ReadFrom(r io.Reader) (c int64, err error) {
	// TODO: We don't know how to read null strings!
	c, s.string, err = readFrom(r)
	return c, err
}

// WriteTo writes the information on SimpleString and dumps it into the Writer
func (s *SimpleString) WriteTo(w io.Writer) (int64, error) {
	if s.isNull {
		return fprintf(w, "%c-1\r\n", byte(SimpleStringOp))
	}
	return fprintf(w, "%c%s\r\n", byte(SimpleStringOp), s.string)
}

// String returns the String representation of the object
func (s *SimpleString) String() string { return s.string }

// Bytes returns the String representation encoded in []bytes
func (s *SimpleString) Bytes() []byte { return []byte(s.string) }
