// Package expire takes care of keeping track of keys that need to expire at some point
package expire

import (
	"container/heap"
	"sync"
	"time"
)

// Expire keeps tracks of keys that must be expired
type Expire struct {
	priorityQueue  *priorityQueue
	mapKeyPosition map[string]*item

	mux sync.Mutex
}

// NewExpire returns a Expire
func NewExpire() *Expire {
	pq := make(priorityQueue, 0)
	return &Expire{
		priorityQueue:  &pq,
		mapKeyPosition: make(map[string]*item),
	}
}

// AddUpdate adds or updates the TTL for a given Key
func (e *Expire) AddUpdate(database int, key string, time int64) {
	e.mux.Lock()
	defer e.mux.Unlock()

	if i, ok := e.mapKeyPosition[key]; ok {
		// Key exists, we update it then.
		e.priorityQueue.update(i, key, database, time)
	} else {
		// Key does not exist. Let's create a new one, then.
		i := &item{database: database, key: key, priority: time}
		heap.Push(e.priorityQueue, i)
		e.mapKeyPosition[key] = i
	}
}

// GetExpired returns the element that
func (e *Expire) GetExpired(time int64) (database int, key string, found bool) {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.priorityQueue.Len() == 0 {
		return 0, "", false
	}

	i, ok := e.priorityQueue.Peek().(*item)
	if !ok {
		panic("unknown stored type; this should never happen")
	}

	if i.priority > time {
		// Top element on the heap has not expired yet, so nothing to do
		return 0, "", false
	}

	// Element expired! Time to remove it from the queue and map, and return it
	heap.Pop(e.priorityQueue)
	delete(e.mapKeyPosition, i.key)
	return i.database, i.key, true
}

func (e *Expire) TTL(key string) (int, bool) {
	e.mux.Lock()
	defer e.mux.Unlock()

	v, ok := e.mapKeyPosition[key]
	if !ok {
		return 0, false
	}

	expireAt := time.Unix(v.priority, 0)
	ttl := int(expireAt.Sub(time.Now()).Seconds())

	return ttl, true
}
