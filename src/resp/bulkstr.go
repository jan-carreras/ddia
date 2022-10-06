package resp

import (
	"fmt"
	"io"
	"strings"
)

type BulkStr struct {
	strings []string
}

func NewBulkStr(strings []string) BulkStr {
	return BulkStr{strings: strings}
}

func (b *BulkStr) Bytes() []byte { return []byte(b.String()) }

func (b *BulkStr) String() string { return strings.Join(b.strings, " ") }

func (b *BulkStr) Strings() []string {
	return b.strings
}

func (b *BulkStr) WriteTo(w io.Writer) (int64, error) {
	length := len(b.strings)
	if length == 0 {
		length = -1
	}

	count := 0
	n, err := fmt.Fprintf(w, "*%d\r\n", length)
	if err != nil {
		return 0, fmt.Errorf("unable to start message: %w", err)
	}
	count += n

	for _, s := range b.strings {
		n, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s)
		if err != nil {
			return 0, fmt.Errorf("unable to write a word in message: %w", err)
		}
		count += n
	}

	return int64(count), nil
}

func (b *BulkStr) ReadFrom(r io.Reader) (int64, error) {
	n, err := b.readFrom(r)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrParsingError, err)
	}

	return n, nil
}

func (b *BulkStr) readFrom(r io.Reader) (int64, error) {
	if err := checkOperation(r, array); err != nil {
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
		case bulkStringOp:
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
