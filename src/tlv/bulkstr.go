package tlv

import (
	"encoding/binary"
	"errors"
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
		return 0, fmt.Errorf("readLength: %w", err)
	}

	for word := 0; word < arrayLength; word++ {
		var operation byte
		if err := binary.Read(r, binary.BigEndian, &operation); err != nil {
			return 0, fmt.Errorf("unable to read operator: %w", err)
		}

		switch operation {
		case '$':
			s := Str{}
			_, err := s.ReadFrom(r)
			if err != nil {
				return 0, fmt.Errorf("str.ReadFrom: %w", err)
			}

			b.strings = append(b.strings, s.String())
		default:
			return 0, errors.New(fmt.Sprintf("unknown operation %q", string(operation)))
		}
	}

	return 0, nil
}

func readLength(r io.Reader) (int, error) {
	var num byte
	arrayLength := 0
	for {
		err := binary.Read(r, binary.BigEndian, &num)
		if err != nil {
			return 0, err
		}

		if num == '\r' { // Stop parsing, consuming one last character
			if err := ignoreDelimiters(r); err != nil {
				return 0, err
			}
			break
		}

		if num < '0' || num > '9' {
			return 0, fmt.Errorf("length must be [0-9]+, %q instead", num)
		}

		arrayLength = (arrayLength * 10) + int(num-'0')
	}

	return arrayLength, nil
}

// Ignores "\r\n" or "\n", failing otherwise
func ignoreDelimiters(r io.Reader) error {
	var char byte

	// Read first character. It should be either \r or \n
	if err := binary.Read(r, binary.BigEndian, &char); err != nil {
		return err
	}

	if !(char == '\r' || char == '\n') {
		return fmt.Errorf("unexpected caracter %q", string(char))
	}

	// If we read \n, we're ignored all delimiters
	if char == '\n' {
		return nil
	}

	// Otherwise, read \n
	if err := binary.Read(r, binary.BigEndian, &char); err != nil {
		return err
	}

	// And make sure it was he character read, otherwise fail
	if char != '\n' {
		return fmt.Errorf("expected character %q, expecting \n", string(char))
	}

	return nil
}
