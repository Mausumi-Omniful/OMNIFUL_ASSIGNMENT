package httpclient

import "sync"

type SafeChannel[T any] interface {
	Read() <-chan T
	Write(T)
	Close()
	IsClosed() bool
}

func NewSafeChannel[T any]() SafeChannel[T] {
	ch := make(chan T)
	return &safeChannel[T]{
		ch:       ch,
		isClosed: false,
	}
}

type safeChannel[T any] struct {
	ch       chan T
	isClosed bool
	mu       sync.Mutex
}

func (c *safeChannel[T]) Read() <-chan T {
	return c.ch
}

func (c *safeChannel[T]) Write(t T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed {
		return
	}
	c.ch <- t
}

func (c *safeChannel[T]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isClosed = true
	close(c.ch)
}

func (c *safeChannel[T]) IsClosed() bool {
	return c.isClosed
}
