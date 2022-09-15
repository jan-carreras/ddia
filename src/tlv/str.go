package tlv

import (
	"errors"
	"fmt"
	"io"
)

type Str struct {
	s string
}

func (b *Str) Bytes() []byte  { return []byte(b.s) }
func (b *Str) String() string { return b.s }

func (b *Str) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (b *Str) ReadFrom(r io.Reader) (int64, error) {
	strLen, err := readLength(r)
	if err != nil {
		return 0, fmt.Errorf("readLength: %w", err)
	}

	buf := make([]byte, strLen)
	written, err := r.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("r.Read(len=%d) string: %w", strLen, err)
	}
	if written != strLen {
		return 0, errors.New(fmt.Sprintf("insufficient data read: expecting Str of length %d, having %d", written, strLen))
	}

	b.s = string(buf)

	// Ignore \r\n
	if err := ignoreDelimiters(r); err != nil {
		return 0, err
	}

	return int64(len(b.s)), nil
}
