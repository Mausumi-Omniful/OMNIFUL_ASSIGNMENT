# Package dmutex

## Overview
The dmutex package provides a distributed mutex implementation using Redis for distributed locking in Go applications. It enables coordination between multiple processes or services by providing thread-safe locking mechanisms across distributed systems.

## Features
- **Distributed Locking**: Safe locking across multiple processes/services
- **Redis-backed**: Uses Redis for lock management
- **Flexible TTL**: Configurable lock timeouts
- **Multiple Lock Modes**: Supports both blocking and non-blocking lock attempts
- **Lock Status Monitoring**: Ability to check lock status and wait for unlocks
- **Context Support**: Full integration with Go's context for cancellation and timeouts

## Key Components

### DMutex Interface
```go
type DMutex interface {
    Lock(ctx context.Context) error
    TryLock(ctx context.Context) (acquired bool, err error)
    Unlock(ctx context.Context) error
    WaitUntilUnlocked(ctx context.Context) (isTTLExpired bool, err error)
    LockedUntil(ctx context.Context) (time.Duration, error)
}
```

## Usage Examples

### Basic Lock/Unlock
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/dmutex"
    "github.com/omniful/go_commons/redis"
    "time"
)

func main() {
    // Initialize Redis client
    redisClient := redis.NewClient(&redis.Config{
        Hosts: []string{"localhost:6379"},
    })

    // Create a new distributed mutex
    mutex := dmutex.New(
        "my-resource-key",
        time.Minute,  // Lock TTL
        redisClient,
    )

    ctx := context.Background()

    // Acquire lock
    if err := mutex.Lock(ctx); err != nil {
        panic(err)
    }

    // ... perform work ...

    // Release lock
    if err := mutex.Unlock(ctx); err != nil {
        panic(err)
    }
}
```

### Non-blocking Lock Attempt
```go
func tryAcquireLock(mutex dmutex.DMutex) {
    ctx := context.Background()
    
    // Try to acquire lock without blocking
    acquired, err := mutex.TryLock(ctx)
    if err != nil {
        panic(err)
    }

    if acquired {
        defer mutex.Unlock(ctx)
        // ... perform work ...
    } else {
        // Lock is held by another process
        fmt.Println("Lock is currently held")
    }
}
```

### Wait for Lock Release
```go
func waitForLock(mutex dmutex.DMutex) {
    ctx := context.Background()
    
    fmt.Println("Waiting for lock to be released...")
    
    isTTLExpired, err := mutex.WaitUntilUnlocked(ctx)
    if err != nil {
        panic(err)
    }

    if isTTLExpired {
        fmt.Println("Lock expired")
    } else {
        fmt.Println("Lock was released")
    }
}
```

### Check Lock Duration
```go
func checkLockDuration(mutex dmutex.DMutex) {
    ctx := context.Background()
    
    duration, err := mutex.LockedUntil(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Lock will expire in: %v\n", duration)
}
```

## Best Practices
1. **Always Use Context**: Provide appropriate context with timeouts for lock operations
2. **Handle Lock Failures**: Implement proper error handling for lock acquisition failures
3. **Use Appropriate TTLs**: Set reasonable lock timeouts based on your use case
4. **Release Locks**: Always ensure locks are released using `defer mutex.Unlock(ctx)`
5. **Check Lock Status**: Use `TryLock` when immediate feedback is needed
6. **Monitor Lock Duration**: Use `LockedUntil` to check remaining lock time

## Error Handling
The package provides specific error types:
- `ErrFailedToLock`: Indicates failure in acquiring the lock
- `ErrFailedToUnlock`: Indicates failure in releasing the lock

```go
if err := mutex.Lock(ctx); err != nil {
    var lockErr *dmutex.ErrFailedToLock
    if errors.As(err, &lockErr) {
        // Handle lock acquisition failure
    }
}
```

## Notes
- The package uses [redsync](https://github.com/go-redsync/redsync) for Redis-based distributed locking
- Lock TTL should be set considering the maximum time your operation might take
- The package is designed to be thread-safe and handle concurrent lock attempts
- Uses Redis pub/sub for efficient lock release notification
- Supports automatic lock expiration through Redis TTL mechanism

## Dependencies
- github.com/go-redis/redis/v8
- github.com/go-redsync/redsync/v4
- github.com/omniful/go_commons/redis
- github.com/omniful/go_commons/dchannel
- github.com/omniful/go_commons/log

## Thread Safety
The implementation is safe for concurrent use by multiple goroutines and processes. 