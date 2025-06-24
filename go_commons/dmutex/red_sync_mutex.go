package dmutex

import (
	"context"
	"errors"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/omniful/go_commons/dchannel"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
	"time"
)

func NewRedSyncMutex(key string, ttl time.Duration, redis *redis.Client) DMutex {
	rp := redsyncredis.NewPool(redis.Cmdable.(goredis.UniversalClient))
	mx := redsync.New(rp).NewMutex(key, redsync.WithExpiry(ttl))
	return &RedSyncMutex{key: key, ttl: ttl, redis: redis, mx: mx}
}

// Use redsync to implement distributed lock.
// https://redis.io/docs/latest/develop/use/patterns/distributed-locks/
type RedSyncMutex struct {
	key string
	ttl time.Duration

	redis *redis.Client
	mx    *redsync.Mutex
	dch   dchannel.DChannel
}

func (m *RedSyncMutex) Lock(ctx context.Context) error {
	if err := m.mx.LockContext(ctx); err != nil {
		return NewErrFailedToLock(err)
	}
	return nil
}

func (m *RedSyncMutex) TryLock(ctx context.Context) (acquired bool, err error) {
	err = m.mx.TryLockContext(ctx)
	if err != nil {
		var errTaken *redsync.ErrTaken
		if errors.As(err, &errTaken) {
			return false, nil
		}
		return false, NewErrFailedToLock(err)
	}
	return true, nil
}

func (m *RedSyncMutex) Unlock(ctx context.Context) error {
	ok, err := m.mx.UnlockContext(ctx)
	if err != nil {
		return NewErrFailedToUnlock(err)
	}
	if !ok {
		return NewErrFailedToUnlock()
	}

	m.getDChannel(true).Close()

	return nil
}

func (m *RedSyncMutex) WaitUntilUnlocked(ctx context.Context) (isTTLExpired bool, err error) {
	ttl, err := m.LockedUntil(ctx)
	if err != nil {
		return false, err
	}

	ch := make(chan bool)
	errCh := make(chan error)
	defer close(ch)
	defer close(errCh)

	// Add a timeout to the context
	ctx, cf := context.WithTimeout(ctx, m.ttl)
	defer cf() // Ensure context gets cancelled while returning

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered from panic: %v", r)
			}
		}()

		dch := m.getDChannel(false)
		defer dch.Close()

		_, err := dch.Listen(ctx)
		if err != nil && !errors.Is(err, dchannel.ErrChanClosed) {
			errCh <- err
			return
		}
		ch <- true
	}()

	select {
	case <-ch:
		return false, nil
	case <-errCh:
		return false, err
	case <-time.After(ttl):
		return true, nil
	}
}

func (m *RedSyncMutex) LockedUntil(ctx context.Context) (time.Duration, error) {
	return m.redis.TTL(ctx, m.key)
}

func (m *RedSyncMutex) getDChannel(refreshChannel bool) dchannel.DChannel {
	if m.dch == nil {
		m.dch = dchannel.New(_generateChannelKey(m.key), m.redis)
	}
	if refreshChannel && m.dch.IsClosed() {
		m.dch = dchannel.New(_generateChannelKey(m.key), m.redis)
	}
	return m.dch
}

func _generateChannelKey(key string) string {
	return fmt.Sprintf("dmutex-%s", key)
}
