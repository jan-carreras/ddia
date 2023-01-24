package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed commands.json
var _commands string

type cmd struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Operation string `json:"operation"`
}

var commands []cmd

func init() {
	err := json.Unmarshal([]byte(_commands), &commands)
	if err != nil {
		panic(fmt.Errorf("unable to decode commands.json: %w", err))
	}
	_commands = "" // No use for this variable once unmarshalled
}

func command(name string) (cmd, bool) {
	name = strings.ToLower(name)
	for _, c := range commands {
		if c.Name == name {
			return c, true
		}
	}
	return cmd{}, false
}

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
	// MGet command
	MGet = "MGET"
	// Set command
	Set = "SET"
	// SetNX command
	SetNX = "SETNX"
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
