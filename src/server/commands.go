package server

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

	/* SERVER */

	// DBSize command
	DBSize = "DBSIZE"

	// FlushDB command
	FlushDB = "FLUSHDB"
	// FlushAll command
	FlushAll = "FLUSHALL"
	// Config command
	Config = "CONFIG"

	/* GENERIC */

	// Exists command
	Exists = "EXISTS"
)
