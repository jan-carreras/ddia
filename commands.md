# 1.0.0
	CONNECTION
		✅ auth: Authenticate to the server
		✅ echo: Echo the given string
		✅ ping: Ping the server
		🏗️  quit: Close the connection
		✅ select: Change the selected database for the current connection
	SERVER
		   bgrewriteaof: Asynchronously rewrite the append-only file
		   bgsave: Asynchronously save the dataset to disk
		✅ dbsize: Return the number of keys in the selected database
		   debug: A container for debugging commands
		✅ flushall: Remove all keys from all databases
		✅ flushdb: Remove all keys from the current database
		   info: Get information and statistics about the server
		   lastsave: Get the UNIX time stamp of the last successful save to disk
		   save: Synchronously save the dataset to disk
		   shutdown: Synchronously save the dataset to disk and then shut down the server
		   slaveof: Make the server a replica of another instance, or promote it as master.
	STRING
		✅ decr: Decrement the integer value of a key by one
		✅ decrby: Decrement the integer value of a key by the given number
		✅ get: Get the value of a key
		🚫 getset: Set the string value of a key and return its old value
		✅ incr: Increment the integer value of a key by one
		✅ incrby: Increment the integer value of a key by the given amount
		✅ mget: Get the values of all the given keys
		✅ set: Set the string value of a key
		✅ setnx: Set the value of a key, only if the key does not exist
		✅ substr: Get a substring of the string stored at a key
	GENERIC
		✅ del: Delete a key
		✅ exists: Determine if a key exists
		   expire: Set a key's time to live in seconds
		   keys: Find all keys matching the given pattern
		   move: Move a key to another database
		✅ randomkey: Return a random key from the keyspace
		✅ rename: Rename a key
		   renamenx: Rename a key, only if the new key does not exist
		   sort: Sort the elements in a list, set or sorted set
		   ttl: Get the time to live for a key in seconds
		   type: Determine the type stored at key
	LIST
		✅ lindex: Get an element from a list by its index
		✅ llen: Get the length of a list
		✅ lpop: Remove and get the first elements in a list
		✅ lpush: Prepend one or multiple elements to a list
		✅ lrange: Get a range of elements from a list
		✅ lrem: Remove elements from a list
		✅ lset: Set the value of an element in a list by its index
		✅ ltrim: Trim a list to the specified range
		✅ rpop: Remove and get the last elements in a list
		✅ rpush: Append one or multiple elements to a list
	SET
		   sadd: Add one or more members to a set
		   scard: Get the number of members in a set
		   sdiff: Subtract multiple sets
		   sdiffstore: Subtract multiple sets and store the resulting set in a key
		   sinter: Intersect multiple sets
		   sinterstore: Intersect multiple sets and store the resulting set in a key
		   sismember: Determine if a given value is a member of a set
		   smembers: Get all the members in a set
		   smove: Move a member from one set to another
		   spop: Remove and return one or multiple random members from a set
		   srandmember: Get one or multiple random members from a set
		   srem: Remove one or more members from a set
		   sunion: Add multiple sets
		   sunionstore: Add multiple sets and store the resulting set in a key
