package resp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var _ dataType = (*Array)(nil)

// Array : Clients send commands to the Redis server using RESP Arrays. Similarly,
// certain Redis commands, that return collections of elements to the client, use
// RESP Arrays as their replies
// Example: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
type Array struct {
	strings []string
}

// NewArray returns an Array type
func NewArray(strings []string) *Array {
	return &Array{strings: strings}
}

// Bytes returns the []bytes representation of an Array
func (b *Array) Bytes() []byte { return []byte(b.String()) }

// String returns the string representation of an Array
func (b *Array) String() string { return strings.Join(b.strings, " ") }

// Strings returns a slice of strings of an Array
func (b *Array) Strings() []string {
	return b.strings
}

// WriteTo writes the array into the Writer. It matches io.WriterTo interface
func (b *Array) WriteTo(w io.Writer) (int64, error) {
	length := len(b.strings)
	if length == 0 {
		length = -1
	}

	buf := bufio.NewWriter(w)

	count := 0
	n, err := fmt.Fprintf(buf, "%c%d\r\n", byte(ArrayOp), length)
	if err != nil {
		return int64(n), fmt.Errorf("unable to start message: %w", err)
	}
	count += n

	for _, s := range b.strings {
		s := s
		n, err := fmt.Fprintf(buf, "%c%d\r\n%s\r\n", byte(BulkStringOp), len(s), s)
		if err != nil {
			return 0, fmt.Errorf("unable to write a word in message: %w", err)
		}
		count += n
	}

	if err := buf.Flush(); err != nil {
		return int64(count), err
	}

	return int64(count), nil
}

// ReadFrom reads an Array object from r. It matches io.ReaderFrom interface
func (b *Array) ReadFrom(r io.Reader) (int64, error) {
	return b.readFrom(r)
}

func (b *Array) readFrom(r io.Reader) (n int64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %v", ErrParsingError, err)
		}
	}()

	arrayLength, err := readLength(r)
	if err != nil {
		return 0, fmt.Errorf("readLength: %v", err)
	}

	for word := 0; word < arrayLength; word++ {
		operation, err := ReadOperation(r)
		if err != nil {
			return 0, fmt.Errorf("unable to read operator: %v", err)
		}

		switch operation {
		case BulkStringOp:
			s := Str{}
			_, err := s.ReadFrom(r)
			if err != nil {
				return 0, err
			}

			b.strings = append(b.strings, s.String())
		case IntegerOp:
			i := Integer{}
			n, err := i.ReadFrom(r)
			if err != nil {
				return n, fmt.Errorf("int.ReadFrom: %v", err)
			}

			b.strings = append(b.strings, i.String())
		default:
			return 0, fmt.Errorf("unknown operator %q", string(operation))
		}
	}

	// TODO: The reporting on the read characters is broken
	return 0, nil
}
