# 1.0.0
	auth: Authenticate to the server
	bgrewriteaof: Asynchronously rewrite the append-only file
	bgsave: Asynchronously save the dataset to disk
	dbsize: Return the number of keys in the selected database
	debug: A container for debugging commands
	decr: Decrement the integer value of a key by one
	decrby: Decrement the integer value of a key by the given number
	del: Delete a key
	echo: Echo the given string
	exists: Determine if a key exists
	expire: Set a key's time to live in seconds
	flushall: Remove all keys from all databases
	flushdb: Remove all keys from the current database
	get: Get the value of a key
	getset: Set the string value of a key and return its old value
	incr: Increment the integer value of a key by one
	incrby: Increment the integer value of a key by the given amount
	info: Get information and statistics about the server
	keys: Find all keys matching the given pattern
	lastsave: Get the UNIX time stamp of the last successful save to disk
	lindex: Get an element from a list by its index
	llen: Get the length of a list
	lpop: Remove and get the first elements in a list
	lpush: Prepend one or multiple elements to a list
	lrange: Get a range of elements from a list
	lrem: Remove elements from a list
	lset: Set the value of an element in a list by its index
	ltrim: Trim a list to the specified range
	mget: Get the values of all the given keys
	move: Move a key to another database
	ping: Ping the server
	quit: Close the connection
	randomkey: Return a random key from the keyspace
	rename: Rename a key
	renamenx: Rename a key, only if the new key does not exist
	rpop: Remove and get the last elements in a list
	rpush: Append one or multiple elements to a list
	sadd: Add one or more members to a set
	save: Synchronously save the dataset to disk
	scard: Get the number of members in a set
	sdiff: Subtract multiple sets
	sdiffstore: Subtract multiple sets and store the resulting set in a key
	select: Change the selected database for the current connection
	set: Set the string value of a key
	setnx: Set the value of a key, only if the key does not exist
	shutdown: Synchronously save the dataset to disk and then shut down the server
	sinter: Intersect multiple sets
	sinterstore: Intersect multiple sets and store the resulting set in a key
	sismember: Determine if a given value is a member of a set
	slaveof: Make the server a replica of another instance, or promote it as master.
	smembers: Get all the members in a set
	smove: Move a member from one set to another
	sort: Sort the elements in a list, set or sorted set
	spop: Remove and return one or multiple random members from a set
	srandmember: Get one or multiple random members from a set
	srem: Remove one or more members from a set
	substr: Get a substring of the string stored at a key
	sunion: Add multiple sets
	sunionstore: Add multiple sets and store the resulting set in a key
	ttl: Get the time to live for a key in seconds
	type: Determine the type stored at key
# 1.0.1
	mset: Set multiple keys to multiple values
	msetnx: Set multiple keys to multiple values, only if none of the keys exist
# 1.0.5
	zrangebyscore: Return a range of members in a sorted set, by score
# 1.2.0
	exec: Execute all commands issued after MULTI
	expireat: Set the expiration for a key as a UNIX timestamp
	multi: Mark the start of a transaction block
	rpoplpush: Remove the last element in a list, prepend it to another list and return it
	zadd: Add one or more members to a sorted set, or update its score if it already exists
	zcard: Get the number of members in a sorted set
	zincrby: Increment the score of a member in a sorted set
	zrange: Return a range of members in a sorted set
	zrem: Remove one or more members from a sorted set
	zremrangebyscore: Remove all members in a sorted set within the given scores
	zrevrange: Return a range of members in a sorted set, by index, with scores ordered from high to low
	zscore: Get the score associated with the given member in a sorted set
# 2.0.0
	append: Append a value to a key
	blpop: Remove and get the first element in a list, or block until one is available
	brpop: Remove and get the last element in a list, or block until one is available
	config: A container for server configuration commands
	discard: Discard all commands issued after MULTI
	get: Get the values of configuration parameters
	hdel: Delete one or more hash fields
	hexists: Determine if a hash field exists
	hget: Get the value of a hash field
	hgetall: Get all the fields and values in a hash
	hincrby: Increment the integer value of a hash field by the given number
	hkeys: Get all the fields in a hash
	hlen: Get the number of fields in a hash
	hmget: Get the values of all the given hash fields
	hmset: Set multiple hash fields to multiple values
	hset: Set the string value of a hash field
	hsetnx: Set the value of a hash field, only if the field does not exist
	hvals: Get all the values in a hash
	psubscribe: Listen for messages published to channels matching the given patterns
	publish: Post a message to a channel
	punsubscribe: Stop listening for messages posted to channels matching the given patterns
	resetstat: Reset the stats returned by INFO
	set: Set configuration parameters to the given values
	setex: Set the value and expiration of a key
	subscribe: Listen for messages published to the given channels
	unsubscribe: Stop listening for messages posted to the given channels
	zcount: Count the members in a sorted set with scores within the given values
	zinterstore: Intersect multiple sorted sets and store the resulting sorted set in a new key
	zrank: Determine the index of a member in a sorted set
	zremrangebyrank: Remove all members in a sorted set within the given indexes
	zrevrank: Determine the index of a member in a sorted set, with scores ordered from high to low
	zunionstore: Add multiple sorted sets and store the resulting sorted set in a new key
# 2.2.0
	brpoplpush: Pop an element from a list, push it to another list and return it; or block until one is available
	getbit: Returns the bit value at offset in the string value stored at key
	linsert: Insert an element before or after another element in a list
	lpushx: Prepend an element to a list, only if the list exists
	persist: Remove the expiration from a key
	rpushx: Append an element to a list, only if the list exists
	setbit: Sets or clears the bit at offset in the string value stored at key
	setrange: Overwrite part of a string at key starting at the specified offset
	strlen: Get the length of the value stored in a key
	unwatch: Forget about all watched keys
	watch: Watch the given keys to determine execution of the MULTI/EXEC block
	zrevrangebyscore: Return a range of members in a sorted set, by score, with scores ordered from high to low
# 2.2.12
	get: Get the slow log's entries
	len: Get the slow log's length
	reset: Clear all entries from the slow log
	slowlog: A container for slow log commands
# 2.2.3
	encoding: Inspect the internal encoding of a Redis object
	idletime: Get the time since a Redis object was last accessed
	object: A container for object introspection commands
	refcount: Get the number of references to the value of the key
# 2.4.0
	client: A container for client connection commands
	getrange: Get a substring of the string stored at a key
	kill: Kill the connection of a client
	list: Get the list of client connections
# 2.6.0
	bitcount: Count set bits in a string
	bitop: Perform bitwise operations between strings
	dump: Return a serialized version of the value stored at the specified key.
	eval: Execute a Lua script server side
	evalsha: Execute a Lua script server side
	exists: Check existence of scripts in the script cache.
	flush: Remove all the scripts from the script cache.
	hincrbyfloat: Increment the float value of a hash field by the given amount
	incrbyfloat: Increment the float value of a key by the given amount
	kill: Kill the script currently in execution.
	load: Load the specified Lua script into the script cache.
	migrate: Atomically transfer a key from a Redis instance to another one.
	pexpire: Set a key's time to live in milliseconds
	pexpireat: Set the expiration for a key as a UNIX timestamp specified in milliseconds
	psetex: Set the value and expiration in milliseconds of a key
	pttl: Get the time to live for a key in milliseconds
	restore: Create a key using the provided serialized value, previously obtained using DUMP.
	script: A container for Lua scripts management commands
	time: Return the current server time
# 2.6.9
	getname: Get the current connection name
	setname: Set the current connection name
# 2.8.0
	channels: List active channels
	hscan: Incrementally iterate hash fields and associated values
	numpat: Get the count of unique patterns pattern subscriptions
	numsub: Get the count of subscribers for channels
	pubsub: A container for Pub/Sub commands
	rewrite: Rewrite the configuration file with the in memory configuration
	scan: Incrementally iterate the keys space
	slaves: List the monitored slaves
	sscan: Incrementally iterate Set elements
	zscan: Incrementally iterate sorted sets elements and associated scores
# 2.8.12
	role: Return the role of the instance in the context of replication
# 2.8.13
	command: Get array of Redis command details
	count: Get total number of Redis commands
	doctor: Return a human readable latency analysis report.
	getkeys: Extract keys given a full Redis command
	graph: Return a latency graph for the event.
	help: Show helpful text about the different subcommands.
	history: Return timestamp-latency samples for the event.
	info: Get array of specific Redis command details, or all when no argument is given.
	latency: A container for latency diagnostics commands
	latest: Return the latest latency samples for all events.
	reset: Reset latency data for one or more events.
# 2.8.4
	flushconfig: Rewrite configuration file
	master: Shows the state of a master
	masters: List the monitored masters
	monitor: Start monitoring
	remove: Stop monitoring
	reset: Reset masters by name pattern
	sentinel: A container for Sentinel commands
	sentinels: List the Sentinel instances
	set: Change the configuration of a monitored master
# 2.8.7
	bitpos: Find first bit set or clear in a string
# 2.8.9
	pfadd: Adds the specified elements to the specified HyperLogLog.
	pfcount: Return the approximated cardinality of the set(s) observed by the HyperLogLog at key(s).
	pfdebug: Internal commands for debugging HyperLogLog values
	pfmerge: Merge N different HyperLogLogs into a single one.
	pfselftest: An internal command for testing HyperLogLog values
	zlexcount: Count the number of members in a sorted set between a given lexicographical range
	zrangebylex: Return a range of members in a sorted set, by lexicographical range
	zremrangebylex: Remove all members in a sorted set between the given lexicographical range
	zrevrangebylex: Return a range of members in a sorted set, by lexicographical range, ordered from higher to lower strings.
# 3.0.0
	addslots: Assign new hash slots to receiving node
	asking: Sent by cluster clients after an -ASK redirect
	bumpepoch: Advance the cluster config epoch
	cluster: A container for cluster commands
	countkeysinslot: Return the number of local keys in the specified hash slot
	delslots: Set hash slots as unbound in receiving node
	failover: Forces a replica to perform a manual failover of its master.
	flushslots: Delete a node's own slots information
	forget: Remove a node from the nodes table
	getkeysinslot: Return local key names in the specified hash slot
	info: Provides info about Redis Cluster node state
	keyslot: Returns the hash slot of the specified key
	meet: Force a node cluster to handshake with another node
	myid: Return the node id
	nodes: Get Cluster config for the node
	pause: Stop processing commands from clients for some time
	readonly: Enables read queries for a connection to a cluster replica node
	readwrite: Disables read queries for a connection to a cluster replica node
	replconf: An internal command for configuring the replication stream
	replicate: Reconfigure a node as a replica of the specified master node
	reset: Reset a Redis Cluster node
	saveconfig: Forces the node to save cluster state on disk
	setslot: Bind a hash slot to a specific node
	slaves: List replica nodes of the specified master node
	slots: Get array of Cluster slot to node mappings
	wait: Wait for the synchronous replication of all the write commands sent in the context of the current connection
# 3.2.0
	bitfield: Perform arbitrary bitfield integer operations on strings
	debug: Set the debug mode for executed scripts.
	geoadd: Add one or more geospatial items in the geospatial index represented using a sorted set
	geodist: Returns the distance between two members of a geospatial index
	geohash: Returns members of a geospatial index as standard geohash strings
	geopos: Returns longitude and latitude of members of a geospatial index
	georadius: Query a sorted set representing a geospatial index to fetch members matching a given maximum distance from a point
	georadiusbymember: Query a sorted set representing a geospatial index to fetch members matching a given maximum distance from a member
	hstrlen: Get the length of the value of a hash field
	reply: Instruct the server whether to reply to commands
# 3.2.1
	touch: Alters the last access time of a key(s). Returns the number of existing keys specified.
# 4.0.0
	doctor: Outputs memory problems report
	freq: Get the logarithmic access frequency counter of a Redis object
	help: Show helpful text about the different subcommands
	list: List all modules loaded by the server
	load: Load a module
	memory: A container for memory diagnostics commands
	module: A container for module commands
	purge: Ask the allocator to release memory
	stats: Show memory usage details
	swapdb: Swaps two Redis databases
	unlink: Delete a key asynchronously in another thread. Otherwise it is just as DEL, but non blocking.
	unload: Unload a module
	usage: Estimate the memory usage of a key
# 5.0.0
	bzpopmax: Remove and return the member with the highest score from one or more sorted sets, or block until one is available
	bzpopmin: Remove and return the member with the lowest score from one or more sorted sets, or block until one is available
	consumers: List the consumers in a consumer group
	create: Create a consumer group.
	delconsumer: Delete a consumer from a consumer group.
	destroy: Destroy a consumer group.
	groups: List the consumer groups of a stream
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	id: Returns the client ID for the current connection
	replicaof: Make the server a replica of another instance, or promote it as master.
	replicas: List replica nodes of the specified master node
	replicas: List the monitored replicas
	setid: Set a consumer group to an arbitrary last delivered ID value.
	stream: Get information about a stream
	unblock: Unblock a client blocked in a blocking command from a different connection
	xack: Marks a pending message as correctly processed, effectively removing it from the pending entries list of the consumer group. Return value of the command is the number of messages successfully acknowledged, that is, the IDs we were actually able to resolve in the PEL.
	xadd: Appends a new entry to a stream
	xclaim: Changes (or acquires) ownership of a message in a consumer group, as if the message was delivered to the specified consumer.
	xdel: Removes the specified entries from the stream. Returns the number of items actually deleted, that may be different from the number of IDs passed in case certain IDs do not exist.
	xgroup: A container for consumer groups commands
	xinfo: A container for stream introspection commands
	xlen: Return the number of entries in a stream
	xpending: Return information and entries from a stream consumer group pending entries list, that are messages fetched but never acknowledged.
	xrange: Return a range of elements in a stream, with IDs matching the specified IDs interval
	xread: Return never seen elements in multiple streams, with IDs greater than the ones reported by the caller for each stream. Can block.
	xreadgroup: Return new entries from a stream using a consumer group, or access the history of the pending entries for a given consumer. Can block.
	xrevrange: Return a range of elements in a stream, with IDs matching the specified IDs interval, in reverse order (from greater to smaller IDs) compared to XRANGE
	xsetid: An internal command for replicating stream values
	xtrim: Trims the stream to (approximately if '~' is passed) a certain size
	zpopmax: Remove and return members with the highest scores in a sorted set
	zpopmin: Remove and return members with the lowest scores in a sorted set
# 6.0.0
	acl: A container for Access List Control commands 
	caching: Instruct the server about tracking or not keys in the next request
	cat: List the ACL categories or the commands inside a category
	deluser: Remove the specified ACL users and the associated rules
	genpass: Generate a pseudorandom secure password to use for ACL users
	getredir: Get tracking notifications redirection client ID if any
	getuser: Get the rules for a specific ACL user
	hello: Handshake with Redis
	help: Show helpful text about the different subcommands
	list: List the current ACL rules in ACL config file format
	load: Reload the ACLs from the configured ACL file
	log: List latest events denied because of ACLs in place
	save: Save the current ACL rules in the configured ACL file
	setuser: Modify or create the rules for a specific ACL user
	tracking: Enable or disable server assisted client side caching support
	users: List the username of all the configured ACL rules
	whoami: Return the name of the user associated to the current connection
# 6.0.6
	lpos: Return the index of matching elements on a list
# 6.2.0
	blmove: Pop an element from a list, push it to another list and return it; or block until one is available
	config: Configure Sentinel
	copy: Copy a key
	createconsumer: Create a consumer in a consumer group.
	failover: Start a coordinated failover between this server and one of its replicas.
	geosearch: Query a sorted set representing a geospatial index to fetch members inside an area of a box or a circle.
	geosearchstore: Query a sorted set representing a geospatial index to fetch members inside an area of a box or a circle, and store the result in another key.
	getdel: Get the value of a key and delete the key
	getex: Get the value of a key and optionally set its expiration
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	help: Show helpful text about the different subcommands
	hrandfield: Get one or multiple random fields from a hash
	info: Returns information about the current client connection.
	lmove: Pop an element from a list, push it to another list and return it
	myid: Get the Sentinel instance ID
	reset: Reset the connection
	smismember: Returns the membership associated with the given elements for a set
	trackinginfo: Return information about server assisted client side caching for the current connection
	unpause: Resume processing of clients that were paused
	xautoclaim: Changes (or acquires) ownership of messages in a consumer group, as if the messages were delivered to the specified consumer.
	zdiff: Subtract multiple sorted sets
	zdiffstore: Subtract multiple sorted sets and store the resulting sorted set in a new key
	zinter: Intersect multiple sorted sets
	zmscore: Get the score associated with the given members in a sorted set
	zrandmember: Get one or multiple random elements from a sorted set
	zrangestore: Store a range of members from sorted set into another key
	zunion: Add multiple sorted sets
# 7.0.0
	addslotsrange: Assign new hash slots to receiving node
	blmpop: Pop elements from a list, or block until one is available
	bzmpop: Remove and return members with scores in a sorted set or block until one is available
	debug: List or update the current configurable parameters
	delete: Delete a function by name
	delslotsrange: Set hash slots as unbound in receiving node
	docs: Get array of specific Redis command documentation
	dryrun: Returns whether the user can execute the given command without executing the command.
	dump: Dump all functions into a serialized binary payload
	expiretime: Get the expiration Unix timestamp for a key
	fcall: Invoke a function
	flush: Deleting all functions
	function: A container for function commands
	getkeysandflags: Extract keys and access flags given a full Redis command
	help: Show helpful text about the different subcommands
	histogram: Return the cumulative distribution of latencies of a subset of commands or all.
	kill: Kill the function currently in execution.
	lcs: Find longest common substring
	links: Returns a list of all TCP links to and from peer nodes in cluster
	list: List information about all the functions
	list: Get an array of Redis command names
	lmpop: Pop elements from a list
	load: Create a function with the given arguments (name, code, description)
	loadex: Load a module with extended parameters
	pexpiretime: Get the expiration Unix timestamp for a key in milliseconds
	restore: Restore all the functions on the given payload
	shardchannels: List active shard channels
	shardnumsub: Get the count of subscribers for shard channels
	shards: Get array of cluster slots to node mappings
	sintercard: Intersect multiple sets and return the cardinality of the result
	spublish: Post a message to a shard channel
	ssubscribe: Listen for messages published to the given shard channels
	stats: Return information about the function currently running (name, description, duration)
	sunsubscribe: Stop listening for messages posted to the given shard channels
	zintercard: Intersect multiple sorted sets and return the cardinality of the result
	zmpop: Remove and return members with scores in a sorted set
# 7.2.0
	myshardid: Return the node shard id
