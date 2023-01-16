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
)

const readBufferSize = 512 // bytes

const (
	/* CONNECTION */

	// Auth command
	Auth = "AUTH"

	// Ping command
	Ping = "PING"
	// Echo command
	Echo = "ECHO"
	// Quit command
	Quit = "QUIT"
	// Select command
	Select = "SELECT"

	/* STRING */

	// Get command
	Get = "GET"
	// Set command
	Set = "SET"
	// Del command
	Del = "DEL"
	// Incr command
	Incr = "INCR"
	// IncrBy command
	IncrBy = "INCRBY"
	// Decr command
	Decr = "DECR"
	// DecrBy command
	DecrBy = "DECRBY"

	/* GENERIC */

	// DBSize command
	DBSize = "DBSIZE"

	// FlushDB command
	FlushDB = "FLUSHDB"
	// FlushAll command
	FlushAll = "FLUSHALL"
)
