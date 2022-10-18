package resp

import (
	"fmt"
	"io"
	"strings"
)

var _ dataType = (*Array)(nil)

// Array are used in order to represent a single binary-safe string up to 512 MB in length.
// They start with $ and
// Example: "$5\r\nhello\r\n"
type Array struct {
	strings []string
}

func NewArray(strings []string) Array {
	return Array{strings: strings}
}

func (b *Array) Bytes() []byte { return []byte(b.String()) }

func (b *Array) String() string { return strings.Join(b.strings, " ") }

func (b *Array) Strings() []string {
	return b.strings
}

func (b *Array) WriteTo(w io.Writer) (int64, error) {
	length := len(b.strings)
	if length == 0 {
		length = -1
	}

	count := 0
	n, err := fmt.Fprintf(w, "%c%d\r\n", byte(ArrayOp), length)
	if err != nil {
		return 0, fmt.Errorf("unable to start message: %w", err)
	}
	count += n

	for _, s := range b.strings {
		n, err := fmt.Fprintf(w, "%c%d\r\n%s\r\n", byte(BulkStringOp), len(s), s)
		if err != nil {
			return 0, fmt.Errorf("unable to write a word in message: %w", err)
		}
		count += n
	}

	return int64(count), nil
}

func (b *Array) ReadFrom(r io.Reader) (int64, error) {
	n, err := b.readFrom(r)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParsingError, err)
	}

	return n, nil
}

func (b *Array) readFrom(r io.Reader) (int64, error) {
	if err := checkOperation(r, ArrayOp); err != nil {
		return 0, fmt.Errorf("checkOperation: %w", err)
	}

	arrayLength, err := readLength(r)
	if err != nil {
		return 0, fmt.Errorf("readLength: %v", err)
	}

	for word := 0; word < arrayLength; word++ {
		r, operation, err := PeakOperation(r)
		if err != nil {
			return 0, fmt.Errorf("unable to read operator: %v", err)
		}

		switch operation {
		case BulkStringOp:
			s := Str{}
			_, err := s.ReadFrom(r)
			if err != nil {
				return 0, fmt.Errorf("str.ReadFrom: %v", err)
			}

			b.strings = append(b.strings, s.String())
		default:
			return 0, fmt.Errorf("unknown operator %q", string(operation))
		}
	}

	// TODO: The reporting on the read characters is broken
	return 0, nil
}
