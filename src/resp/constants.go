package resp

// More documentation on: https://redis.io/docs/reference/protocol-spec/

const (
	// SimpleStringOp defines a simple string datatype. Example: "+OK\r\n"
	SimpleStringOp = '+'
	// ErrorOp defines a error. Example: "-Error message\r\n"
	ErrorOp = '-'
	// IntegerOp defines a error Integer. Example: ":1000\r\n"
	IntegerOp = ':'
	// BulkStringOp defines a BulkStringOp. Example: "$5\r\nhello\r\n"
	BulkStringOp = '$'
	// ArrayOp defines a RESP Arrays. Example: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	//
	// Clients send commands to the Redis server using RESP Arrays. Similarly,
	// certain Redis commands, that return collections of elements to the client, use
	// RESP Arrays as their replies.
	ArrayOp = '*'

	// RawPing it's a made up name. The redis-benchmark utility sends "PING" (instead of "*1\r\n$4\r\nPING\r\n")
	// and I think it's kinda like a heartbeat/healthcheck mechanism
	RawPing = 'P'
)

const readBufferSize = 512 // bytes
