package resp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Decode decodes the network input into the appropriate object. Returns a fmt.Stringer that can return
// the internal representation of the object as string. It's convinient for debugging.
func Decode(input io.Reader) (fmt.Stringer, error) {
	operation, err := readOperation(input)
	if err != nil {
		return nil, err
	}

	var dt DataType
	switch operation {
	case SimpleStringOp:
		dt = &SimpleString{}
	case IntegerOp:
		dt = &Integer{}
	case BulkStringOp:
		dt = &Str{}
	case ArrayOp:
		dt = &Array{}
	case ErrorOp:
		dt = &Error{}
	default:
		return nil, fmt.Errorf("unknown operation type %q", operation)
	}

	_, err = dt.ReadFrom(input)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

// ReadOperation returns the operation type on the stream (eg: +, - , *, ...)
func readOperation(r io.Reader) (byte, error) {
	var operation byte
	err := binary.Read(r, binary.BigEndian, &operation)
	if err != nil {
		return 0, err
	}

	return operation, nil
}
