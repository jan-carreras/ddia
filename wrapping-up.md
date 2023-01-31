# Wrapping Up

I've spent January playing around with this Redis implementation. It has been an excellent exercise, and I've learnt quite a lot.

Things I needed more time to do: Snapshots (RDB), read replicas, high availability, compactation of AOF, being able to store more info than available RAM, and many, many others.

The good news is that I know how I would approach those things now — at least, how to start thinking about them and planning the implementation.

That's what I was aiming to do. I'm sure I'll repeat the exercise in the future (implementation in Rust, anyone?). Probably implementing other parts (Pub/Sub, replication / other commands / ... ).

I would recommend this exercise to anyone wanting to learn a new programming language or how to program a "network service".

# What I've learnt

* **Build from simple principles**
	* Simple, human-readable network protocol goes very, very far
	* Keeping everything simple is key
	* Delaying generic implementations until you know the Domain very well.
* **Design Documents are key** (for me)
	* It helps me organise my thoughts and define the scope 
	* Re-read the design document after two days to spot missing things
	* Tiny things (reading config) can have many edge cases (cyclic imports, glob file expansion, multi-directives for a key, ...)
	* Adding new features will mess with previous designs
* **Beautiful Designs get ruined by "small stuff"** (reading original C Redis code):
	*  Adding many configuration options, observability, logging, error handling, client/server TCP stuff (timeouts, buffers, ...), performance optimisations
	* You have to be very intentional when adding new stuff not to ruin the rest of the stuff
* **Test, test and test**
	* That's a favourite of mine. Not really something I've learnt but it proves worth it every single time
* **Minimise**:
	* Dependencies: stdlib is very good already, use it
	* Interface surface: use stdlib if you can (eg: io.ReadWritter, io.Discard)
	* Code written: maximising readability
	* Complexity: KISS
creating new standards: use uber style standard, or google's one and stick to it

