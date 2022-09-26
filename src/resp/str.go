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
	out := fmt.Sprintf("%c%d\r\n%s\r\n", bulkStringOp, len(b.s), b.s)
	if _, err := w.Write([]byte(out)); err != nil {
		return 0, fmt.Errorf("writing string operator: %w", err)
	}

	return 0, nil
}

func (b *Str) ReadFrom(r io.Reader) (int64, error) {
	if err := checkOperation(r, bulkStringOp); err != nil {
		return 0, fmt.Errorf("checkOperation: %w", err)
	}

	strLen, err := readLength(r)
	if err != nil {
		return 0, fmt.Errorf("readLength: %v: %w", err, ErrParsingError)
	}

	buf := make([]byte, strLen)

	read, err := r.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("r.Read(len=%d) string: %v: %w", strLen, err, ErrParsingError)
	}

	if read != strLen {
		return 0, fmt.Errorf("insufficient data read: expecting Str of length %d, having %d : %w", read, strLen, ErrParsingError)
	}

	b.s = string(buf)

	// Ignore \r\n
	if err := ignoreDelimiters(r); err != nil {
		return 0, fmt.Errorf("%v: %w", err, ErrParsingError)
	}

	return int64(len(b.s)), nil
}
