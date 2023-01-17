package resp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ErrParsingError returned when we cannot parse some information
var ErrParsingError = errors.New("parsing error")

// ErrEncodingError returned when we cannot encode some information
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
		return fmt.Errorf("expected character %q, expecting \\n", string(char))
	}

	return nil
}

// PeakOperation returns the next character on the reader, and a new reader with
// that character in the stream.
func PeakOperation(r io.Reader) (io.Reader, byte, error) {
	var operation byte
	if err := binary.Read(r, binary.BigEndian, &operation); err != nil {
		return r, 0, fmt.Errorf("unable to read operator: %w", err)
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(operation)
	r = io.MultiReader(buf, r) // Load again the read character into the reader stream

	return r, operation, nil
}

func ReadOperation(r io.Reader) (byte, error) {
	return readOperation(r)
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

// ignoreDelimiterCharacters ignores the last two characters if they are \r\n or fails
func ignoreDelimiterCharacters(s string) (string, error) {
	if l := len(s); l < 2 {
		return "", fmt.Errorf("invalid string lenght")
	} else if s[l-2] != '\r' || s[l-1] != '\n' {
		return "", fmt.Errorf("unexpcted end")
	} else {
		s = s[:l-2] // Ignore the last two characters
	}

	return s, nil
}

// readFrom reads all the content of the reader and returns an string
func readFrom(r io.Reader) (readCount int64, s string, err error) {
	buf := make([]byte, readBufferSize)
	for {
		c, err := r.Read(buf)
		readCount += int64(c)

		if errors.Is(err, io.EOF) || c == 0 {
			break
		}

		if err != nil {
			return readCount, s, fmt.Errorf("unable to read: %w", err)
		}

		s += string(buf[:c])

		// We've read all there was to read, we can stop looking for data
		if c < readBufferSize {
			break
		}
	}

	s, err = ignoreDelimiterCharacters(s)
	if err != nil {
		return readCount, "", fmt.Errorf("ignoreDelimiterCharacters: %w", err)
	}

	return readCount, s, nil
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadSlice('\n')
	if err != nil {
		if !errors.Is(err, bufio.ErrBufferFull) {
			return nil, err
		}

		l := make([]byte, len(line))
		copy(l, line)
		line, err = reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		l = append(l, line...)
		line = l
	}

	if len(line) < 2 || line[len(line)-1] != '\n' || line[len(line)-2] != '\r' {
		return nil, errors.New("invalid response")
	}

	return line[:len(line)-2], err
}

func fprintf(w io.Writer, format string, a ...any) (int64, error) {
	n, err := fmt.Fprintf(w, format, a...)
	return int64(n), err
}
