package resp

import (
	"io"
	"strconv"
)

var _ DataType = (*Integer)(nil)

// Integer datatype defines a Integer value
type Integer struct {
	value string
}

// NewInteger returns a Integer datatype
func NewInteger(i int) *Integer {
	return &Integer{value: strconv.Itoa(i)}
}

// ReadFrom reads from the Reader and loads the Integer object
func (s *Integer) ReadFrom(r io.Reader) (c int64, err error) {
	c, s.value, err = readFrom(r)
	return c, err
}

// WriteTo writes the information on Integer and dumps it into the Writer
func (s *Integer) WriteTo(w io.Writer) (int64, error) {
	return fprintf(w, "%c%s\r\n", byte(IntegerOp), s.value)
}

// String returns the String representation of the object
func (s *Integer) String() string { return s.value }

// Bytes returns the String representation encoded in []bytes
func (s *Integer) Bytes() []byte { return []byte(s.value) }
