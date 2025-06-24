// Distributed Mutex
package dmutex

import (
	"context"
	"github.com/omniful/go_commons/redis"
	"time"
)

type DMutex interface {
	// Lock would try to acquire lock and would wait until acquired
	Lock(ctx context.Context) error
	// TryLock would attempt to acquire lock and return false if already locked
	TryLock(ctx context.Context) (acquired bool, err error)
	// Unlock the acquired lock
	Unlock(ctx context.Context) error
	// WaitUntilUnlocked would block until unlocked
	WaitUntilUnlocked(ctx context.Context) (isTTLExpired bool, err error)
	// LockedUntil returns the TTL of the lock
	LockedUntil(ctx context.Context) (time.Duration, error)
}

func New(key string, ttl time.Duration, redis *redis.Client) DMutex {
	return NewRedSyncMutex(key, ttl, redis)
}
