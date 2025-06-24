# Package dchannel

## Overview
The dchannel package provides a distributed channel implementation using Redis pub/sub for communication across distributed systems. It enables reliable message passing and synchronization between different processes or services in a distributed environment.

## Features
- **Distributed Communication**: Message passing across multiple processes/services
- **Redis Pub/Sub**: Built on Redis for reliable message delivery
- **Context Support**: Full integration with Go's context for cancellation and timeouts
- **Channel Lifecycle**: Managed channel creation and cleanup
- **Thread Safety**: Safe for concurrent use
- **Status Monitoring**: Track channel state and health

## Key Components

### DChannel Interface
```go
type DChannel interface {
    // Listen blocks until a message is received or the channel is closed
    Listen(ctx context.Context) (interface{}, error)
    
    // Close closes the channel and notifies all listeners
    Close() error
    
    // CloseLockedUntil closes the channel's lock duration tracking
    CloseLockedUntil() error
    
    // IsClosed returns true if the channel has been closed
    IsClosed() bool
}
```

## Usage Examples

### Basic Channel Usage
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/dchannel"
    "github.com/omniful/go_commons/redis"
)

func main() {
    // Initialize Redis client
    redisClient := redis.NewClient(&redis.Config{
        Hosts: []string{"localhost:6379"},
    })

    // Create a new distributed channel
    channel := dchannel.New("my-channel", redisClient)

    ctx := context.Background()

    // Start listening for messages in a goroutine
    go func() {
        msg, err := channel.Listen(ctx)
        if err != nil {
            panic(err)
        }
        // Process the message
        fmt.Printf("Received message: %v\n", msg)
    }()

    // Close the channel when done
    defer channel.Close()
}
```

### Channel with Timeout
```go
func listenWithTimeout(channel dchannel.DChannel) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Listen for message with timeout
    msg, err := channel.Listen(ctx)
    if err != nil {
        if err == context.DeadlineExceeded {
            fmt.Println("Listen timeout exceeded")
            return
        }
        panic(err)
    }

    fmt.Printf("Received message: %v\n", msg)
}
```

### Check Channel Status
```go
func checkChannelStatus(channel dchannel.DChannel) {
    if channel.IsClosed() {
        fmt.Println("Channel is closed")
        return
    }
    fmt.Println("Channel is open")
}
```

## Best Practices
1. **Always Close Channels**: Use `defer channel.Close()` to ensure proper cleanup
2. **Handle Context Cancellation**: Provide appropriate context with timeouts
3. **Check Channel Status**: Use `IsClosed()` to verify channel state before operations
4. **Error Handling**: Implement proper error handling for channel operations
5. **Resource Management**: Close channels when they're no longer needed
6. **Concurrent Access**: Design for concurrent access when multiple goroutines use the channel

## Error Handling
The package provides specific error types:
- `ErrChanClosed`: Indicates that the channel has been closed
- Other Redis-related errors may be returned during operations

```go
msg, err := channel.Listen(ctx)
if err != nil {
    if errors.Is(err, dchannel.ErrChanClosed) {
        // Handle closed channel
        return
    }
    // Handle other errors
}
```

## Notes
- Uses Redis pub/sub mechanism for message distribution
- Channels are identified by unique keys in Redis
- Channel closure is propagated to all listeners
- Supports automatic cleanup of Redis resources
- Designed for use in distributed systems like the dmutex package

## Dependencies
- github.com/go-redis/redis/v8
- github.com/omniful/go_commons/redis
- github.com/omniful/go_commons/log

## Thread Safety
The implementation is safe for concurrent use by multiple goroutines and processes. All operations are thread-safe and can be used in concurrent environments.

## Performance Considerations
1. **Message Size**: Keep messages small to minimize network overhead
2. **Listener Count**: Be mindful of the number of concurrent listeners
3. **Redis Connection**: Maintain stable Redis connection for reliable operation
4. **Context Usage**: Use appropriate context timeouts to prevent hanging operations 