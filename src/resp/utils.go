package resp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var ErrParsingError = errors.New("parsing error")
var ErrEncodingError = errors.New("encoding error")

// dataType is the interface to be implemented by each datatype on RESP
type dataType interface {
	io.WriterTo
	io.ReaderFrom
	fmt.Stringer
}

func readLength(r io.Reader) (int, error) {
	var num byte
	arrayLength := 0
	for i := 0; ; i++ {
		err := binary.Read(r, binary.BigEndian, &num)
		if err != nil {
			return 0, err
		}

		// The first element should always be numeric
		if i == 0 && (num < '0' || num > '9') {
			return 0, fmt.Errorf("[first] length must be [0-9]+, %q instead", num)
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
		return fmt.Errorf("unexpected character %q", string(char))
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

// PeakOperation returns the next character on the reader, and a new reader with
// that character in the stream.
func PeakOperation(r io.Reader) (io.Reader, byte, error) {
	var operation byte
	if err := binary.Read(r, binary.BigEndian, &operation); err != nil {
		return r, 0, fmt.Errorf("unable to read operator: %v", err)
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(operation)
	r = io.MultiReader(buf, r) // Load again the read character into the reader stream

	return r, operation, nil
}

func readOperation(r io.Reader) (byte, error) {
	var operation byte
	err := binary.Read(r, binary.BigEndian, &operation)
	if err != nil {
		return 0, err
	}

	return operation, nil
}

func checkOperation(r io.Reader, expectedOperation byte) error {
	operation, err := readOperation(r)
	if err != nil {
		return fmt.Errorf("readOperation: %w", err)
	}

	if operation != expectedOperation {
		return fmt.Errorf(
			"unknown operation: expecting %q, have %q",
			expectedOperation,
			operation,
		)
	}

	return nil
}
