package resp

import (
	"fmt"
	"io"
)

var _ dataType = (*Integer)(nil)

type Integer struct {
	string string
}

func NewInteger(string string) *Integer {
	return &Integer{string: string}
}

// ReadFrom reads from the Reader and loads the Integer object
// Example: "+OK\r\n"
func (s *Integer) ReadFrom(r io.Reader) (readCount int64, err error) {
	err = checkOperation(r, IntegerOp)
	readCount += 1 // Read the first byte
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

// WriteTo writes the information on Integer and dumps it into the Writer
func (s *Integer) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c%s\r\n", byte(IntegerOp), s.string)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

// String returns the String representation of the object
func (s *Integer) String() string {
	return s.string
}

// Bytes returns the String representation encoded in []bytes
func (s *Integer) Bytes() []byte {
	return []byte(s.string)
}
