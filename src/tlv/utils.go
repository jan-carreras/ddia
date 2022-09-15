package tlv

import (
	"encoding/binary"
	"fmt"
	"io"
)

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
