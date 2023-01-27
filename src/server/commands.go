package server

import (
	"strings"
)

type cmd struct {
	Name      string `json:"name"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
	Kind      string `json:"kind"`
}

// Commands describe all the commands supported (or not) by the Server.
//
//go:generate go run gen.go
var Commands = []cmd{
	// String commands
	{Name: "Get", Operation: "read", Status: "implemented", Kind: "string"},
	{Name: "Set", Operation: "write", Status: "implemented", Kind: "string"},
	{Name: "GetSet", Operation: "write", Status: "won't-do", Kind: "string"},
	{Name: "Incr", Operation: "write", Status: "implemented", Kind: "string"},
	{Name: "IncrBy", Operation: "write", Status: "implemented", Kind: "string"},
	{Name: "Decr", Operation: "write", Status: "implemented", Kind: "string"},
	{Name: "DecrBy", Operation: "write", Status: "implemented", Kind: "string"},
	{Name: "Substr", Operation: "read", Status: "implemented", Kind: "string"},
	{Name: "MGet", Operation: "read", Status: "implemented", Kind: "string"},
	// Connection commands
	{Name: "Echo", Operation: "read", Status: "implemented", Kind: "connection"},
	{Name: "Ping", Operation: "read", Status: "implemented", Kind: "connection"},
	{Name: "Quit", Operation: "read", Status: "partially-implemented", Kind: "connection"},
	{Name: "Select", Operation: "read", Status: "implemented", Kind: "connection"},
	{Name: "Auth", Operation: "read", Status: "implemented", Kind: "connection"},
	// Generic commands
	{Name: "Del", Operation: "write", Status: "implemented", Kind: "generic"},
	{Name: "Exists", Operation: "read", Status: "implemented", Kind: "generic"},
	{Name: "Move", Operation: "write", Status: "implemented", Kind: "generic"},
	{Name: "RandomKey", Operation: "read", Status: "implemented", Kind: "generic"},
	{Name: "Rename", Operation: "write", Status: "implemented", Kind: "generic"},
	{Name: "Expire", Operation: "write", Status: "implemented", Kind: "generic"},
	{Name: "TTL", Operation: "read", Status: "implemented", Kind: "generic"},
	// Server commands
	{Name: "DBSize", Operation: "read", Status: "implemented", Kind: "server"},
	{Name: "FlushDB", Operation: "write", Status: "implemented", Kind: "server"},
	{Name: "FlushAll", Operation: "write", Status: "implemented", Kind: "server"},
	{Name: "Config", Operation: "write", Status: "partially-implemented", Kind: "server"},
	// List commands
	{Name: "SetNX", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "LLen", Operation: "read", Status: "implemented", Kind: "list"},
	{Name: "LRange", Operation: "read", Status: "implemented", Kind: "list"},
	{Name: "LRem", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "LIndex", Operation: "read", Status: "implemented", Kind: "list"},
	{Name: "LSet", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "LPush", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "RPush", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "LPop", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "RPop", Operation: "write", Status: "implemented", Kind: "list"},
	{Name: "LTrim", Operation: "write", Status: "implemented", Kind: "list"},
}

func getCommand(name string) (cmd, bool) {
	name = strings.ToLower(name)
	for _, c := range Commands {
		if strings.ToLower(c.Name) == name {
			return c, true
		}
	}
	return cmd{}, false
}
