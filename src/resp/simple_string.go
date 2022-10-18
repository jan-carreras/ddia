package resp

import (
	"errors"
	"fmt"
	"io"
)

type SimpleString struct {
	string string
}

func NewSimpleString(string string) *SimpleString {
	return &SimpleString{string: string}
}

// ReadFrom reads from the Reader and loads the SimpleString object
// Example: "+OK\r\n"
func (s *SimpleString) ReadFrom(r io.Reader) (readCount int64, err error) {
	err = checkOperation(r, SimpleStringOp)
	readCount += 1 // Read the first byte
	if err != nil {
		return readCount, err
	}

	c, err := s.readFrom(r)
	readCount += c
	if err != nil {
		return readCount, fmt.Errorf("readFrom: %w", err)
	}

	return readCount, nil
}

func (s *SimpleString) readFrom(r io.Reader) (readCount int64, err error) {
	buf := make([]byte, readBufferSize)
	for {
		c, err := r.Read(buf)
		readCount += int64(c)

		if errors.Is(err, io.EOF) || c == 0 {
			break
		}

		if err != nil {
			return readCount, fmt.Errorf("unable to read: %w", err)
		}

		s.string += string(buf[:c])
	}

	if err := s.ignoreDelimiterCharacters(); err != nil {
		return readCount, fmt.Errorf("ignoreDelimiterCharacters: %w", err)
	}

	return readCount, nil
}

// ignoreDelimiterCharacters ignores the last two characters if they are \r\n or fails
func (s *SimpleString) ignoreDelimiterCharacters() error {
	if l := len(s.string); l < 2 {
		return fmt.Errorf("invalid string lenght")
	} else if s.string[l-2] != '\r' || s.string[l-1] != '\n' {
		fmt.Printf("%q\n\n", s.string)
		fmt.Printf(s.string)
		return fmt.Errorf("unexpcted end")
	} else {
		s.string = s.string[:l-2] // Ignore the last two characters
	}

	return nil
}

// WriteTo writes the information on SimpleString and dumps it into the Writer
func (s *SimpleString) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c%s\r\n", byte(SimpleStringOp), s.string)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

// String returns the String representation of the object
func (s *SimpleString) String() string {
	return s.string
}
