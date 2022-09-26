package resp

import (
	"fmt"
	"io"
	"strings"
)

type BulkStr struct {
	strings []string
}

func (b *BulkStr) Bytes() []byte  { return []byte(b.String()) }
func (b *BulkStr) String() string { return strings.Join(b.strings, " ") }

func (b *BulkStr) WriteTo(_ io.Writer) (int64, error) {
	return 0, nil
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
		r, operation, err := peakOperation(r)
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
