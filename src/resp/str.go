package resp

import (
	"fmt"
	"io"
)

var _ dataType = (*Str)(nil)

// Str returns a Str datatype
type Str struct {
	s      string
	isNull bool
}

// NewStr returns a Str
func NewStr(s string) *Str {
	return &Str{s: s}
}

func NewNullStr() *Str {
	return &Str{isNull: true}
}

// Bytes returns the bytes representation of the string
func (b *Str) Bytes() []byte { return []byte(b.s) }

// String returns the string representation of Str datatype
func (b *Str) String() string { return b.s }

// WriteTo writes the information on Str and dumps it into the Writer
func (b *Str) WriteTo(w io.Writer) (int64, error) {
	if b.isNull {
		return fprintf(w, "$-1\r\n")
	}
	return fprintf(w, "%c%d\r\n%s\r\n", byte(BulkStringOp), len(b.s), b.s)
}

// ReadFrom reads from the Reader and loads the Str object
func (b *Str) ReadFrom(r io.Reader) (int64, error) {
	return b.readFrom(r)
}

func (b *Str) readFrom(r io.Reader) (n int64, err error) {
	// TODO: We cannot parse null strings "$-1\r\n"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %v", ErrParsingError, err)
		}
	}()

	var readCount int64

	strLen, err := readLength(r)
	if err != nil {
		return readCount, fmt.Errorf("readLength: %v", err)
	}

	buf := make([]byte, strLen)

	read := 0
	for read != strLen {
		c, err := r.Read(buf[read:])
		read += c
		readCount += int64(read)
		if err != nil {
			return readCount, fmt.Errorf("r.Read(len=%d): %v", strLen, err)
		}

		if read == 0 {
			break
		}
	}

	if read != strLen {
		return 0, fmt.Errorf("insufficient data read: read %d, want %d", read, strLen)
	}

	b.s = string(buf)

	// Ignore \r\n
	if err := ignoreDelimiters(r); err != nil {
		return readCount, err
	}

	return int64(len(b.s)), nil
}
