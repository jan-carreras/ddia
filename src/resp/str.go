package resp

import (
	"fmt"
	"io"
)

type Str struct {
	s string
}

func (b *Str) Bytes() []byte  { return []byte(b.s) }
func (b *Str) String() string { return b.s }

func (b *Str) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "%c%d\r\n%s\r\n", byte(BulkStringOp), len(b.s), b.s)
	if err != nil {
		return 0, fmt.Errorf("%w: writing string operator: %v", ErrEncodingError, err)
	}

	return int64(n), nil
}

func (b *Str) ReadFrom(r io.Reader) (int64, error) {
	n, err := b.readFrom(r)
	if err != nil {
		return n, fmt.Errorf("%w: %v", ErrParsingError, err)
	}

	return n, nil
}

func (b *Str) readFrom(r io.Reader) (int64, error) {
	var readCount int64

	if err := checkOperation(r, BulkStringOp); err != nil {
		return readCount, fmt.Errorf("checkOperation: %w", err)
	}

	strLen, err := readLength(r)
	if err != nil {
		return readCount, fmt.Errorf("readLength: %v", err)
	}

	buf := make([]byte, strLen)

	read, err := r.Read(buf)
	readCount += int64(read)
	if err != nil {
		return readCount, fmt.Errorf("r.Read(len=%d): %v", strLen, err)
	}

	if read != strLen {
		return 0, fmt.Errorf("insufficient data read: expecting Str of length %d, having %d", read, strLen)
	}

	b.s = string(buf)

	// Ignore \r\n
	if err := ignoreDelimiters(r); err != nil {
		return readCount, err
	}

	return int64(len(b.s)), nil
}
