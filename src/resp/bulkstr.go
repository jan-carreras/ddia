package resp

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type BulkStr struct {
	strings []string
}

func (b *BulkStr) Bytes() []byte  { return []byte(b.String()) }
func (b *BulkStr) String() string { return strings.Join(b.strings, " ") }

func (b *BulkStr) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (b *BulkStr) ReadFrom(r io.Reader) (int64, error) {
	arrayLength, err := readLength(r)
	if err != nil {
		return 0, fmt.Errorf("readLength: %w: %v", ErrParsingError, err)
	}

	for word := 0; word < arrayLength; word++ {
		var operation byte
		if err := binary.Read(r, binary.BigEndian, &operation); err != nil {
			return 0, fmt.Errorf("unable to read operator: %w: %v", ErrParsingError, err)
		}

		switch operation {
		case '$':
			s := Str{}
			_, err := s.ReadFrom(r)
			if err != nil {
				return 0, fmt.Errorf("str.ReadFrom: %w: %v", ErrParsingError, err)
			}

			b.strings = append(b.strings, s.String())
		default:
			return 0, fmt.Errorf("unknown operator %q: %w", string(operation), ErrParsingError)
		}
	}

	// TODO: The reporting on the read characters is broken
	return 0, nil
}
