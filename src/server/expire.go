package server

import (
	"context"
	"time"
)

// lookForKeysToExpire to called as goroutine. Every second it will look for keys to invalidate. It will invalidate
// all keys that are stale. To stop it, close the context.
func (s *Server) lookForKeysToExpire(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	findAndExpire := func() {
		for {
			database, key, isThereSomethingToExpire := s.expire.GetExpired(time.Now().Unix())
			if !isThereSomethingToExpire {
				break
			}

			s.options.dbs[database].Lock()
			s.options.dbs[database].Del(key)
			s.options.dbs[database].Unlock()
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			findAndExpire()
		}
	}
}
