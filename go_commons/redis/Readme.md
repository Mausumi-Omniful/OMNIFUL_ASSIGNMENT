# Redis Package

A comprehensive Go Redis client package that provides a high-level interface for Redis operations, supporting both standalone and cluster modes. The package is built on top of [go-redis/redis](https://github.com/go-redis/redis) with additional features and simplified interfaces.

## Features

- Supports both standalone and cluster mode Redis configurations
- Built-in connection pooling and configuration
- New Relic integration for monitoring
- Comprehensive API for various Redis data types:
  - Strings
  - Sorted Sets
  - Hash Maps
  - Pub/Sub
  - Sets
  - Keys Operations
  - Pipeline Operations

## Installation

```bash
go get github.com/omniful/go_commons/redis
```

## Configuration

The package provides a flexible configuration structure:

```go
type Config struct {
    // Whether redis is running in cluster mode
    ClusterMode bool

    // Read preferences for cluster mode
    ServeReadsFromSlaves bool
    ServeReadsFromMasterAndSlaves bool

    // Connection pool configuration
    PoolSize uint      // Maximum number of connections
    PoolFIFO bool      // true for FIFO pool, false for LIFO pool
    MinIdleConn uint   // Minimum number of idle connections

    // Database selection (non-cluster mode only)
    DB uint

    // Redis server addresses
    Hosts []string

    // Timeouts
    DialTimeout time.Duration
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    IdleTimeout time.Duration
}
```

## Usage Examples

### Initializing the Client

```go
import (
    "github.com/omniful/go_commons/redis"
)

config := &redis.Config{
    Hosts: []string{"localhost:6379"},
    PoolSize: 50,
    MinIdleConn: 10,
}

client := redis.NewClient(config)
defer client.Close()
```

### String Operations

```go
ctx := context.Background()

// Set a key with expiration
success, err := client.Set(ctx, "key", "value", 1*time.Hour)

// Get a key
value, err := client.Get(ctx, "key")

// Set multiple keys
err = client.MSet(ctx, "key1", "value1", "key2", "value2")

// Get multiple keys
values, err := client.MGet(ctx, "key1", "key2")

// Increment/Decrement
newVal, err := client.Incr(ctx, "counter")
newVal, err := client.Decr(ctx, "counter")
```

### Hash Operations

```go
ctx := context.Background()

// Set a single hash field
count, err := client.HSet(ctx, "user:1", "name", "John")

// Set multiple hash fields
fields := map[string]interface{}{
    "name": "John",
    "age":  30,
    "city": "New York",
}
count, err = client.HSetAll(ctx, "user:1", fields)

// Get a single field
name, err := client.HGet(ctx, "user:1", "name")

// Get multiple fields
values, err := client.HMGet(ctx, "user:1", "name", "age")

// Get all fields
allFields, err := client.HGetAll(ctx, "user:1")

// Delete fields
deleted, err := client.HDel(ctx, "user:1", "age", "city")

// Increment numeric fields
newVal, err := client.HIncrBy(ctx, "user:1", "visits", 1)
newFloat, err := client.HIncrByFloat(ctx, "user:1", "score", 1.5)
```

### Sets Operations

```go
ctx := context.Background()

// Add elements to a set
count, err := client.SAdd(ctx, "myset", "element1", "element2", "element3")

// Check if element exists in set
exists, err := client.SIsMember(ctx, "myset", "element1")
```

### Key Operations

```go
ctx := context.Background()

// Check if keys exist
count, err := client.Exists(ctx, "key1", "key2")

// Delete keys (blocking)
deleted, err := client.Del(ctx, "key1", "key2")

// Delete keys (non-blocking/async)
deleted, err := client.Unlink(ctx, "key1", "key2")
```

### Pipeline Operations

Pipeline allows you to send multiple commands to Redis server in a single request, which can significantly improve performance by reducing network round trips.

```go
ctx := context.Background()

// Create a pipeline
pipe := client.Pipeline()

// Queue multiple commands
pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
pipe.Get(ctx, "key1")
pipe.Get(ctx, "key2")

// Execute all commands in a single request
cmds, err := pipe.ExecContext(ctx)
if err != nil {
    // Handle error
}

// Process results
for _, cmd := range cmds {
    if cmd.Err() != nil {
        // Handle command error
    }
    // Process command result
}
```

### Sorted Sets Operations

```go
ctx := context.Background()

// Add members to sorted set
members := []redis.ZMember{
    {Member: "member1", Score: 1.0},
    {Member: "member2", Score: 2.0},
}
count, err := client.ZAdd(ctx, "myset", members...)

// Get range by score
items, err := client.ZRangeByScore(ctx, "myset", "0", "2")

// Remove members
removed, err := client.ZRem(ctx, "myset", "member1", "member2")

// Get score of a member
score, err := client.ZScore(ctx, "myset", "member1")
```

### Publish/Subscribe

```go
ctx := context.Background()

// Publishing messages
count, err := client.Publish(ctx, "channel1", "Hello World")

// Subscribing to channels
pubsub, err := client.Subscribe(ctx, "channel1")
if err != nil {
    // Handle error
}
defer pubsub.Close()

// Using channel-based subscription
ch, err := client.SubscribeChannel(ctx, "channel1")
if err != nil {
    // Handle error
}

// Reading messages
for msg := range ch {
    fmt.Printf("Received message: %s\n", msg.Payload)
}
```

### Counter Operations

```go
ctx := context.Background()

// Rate limiting with IncrWithLimit
// Increments counter with a maximum limit and expiry
newVal, err := client.IncrWithLimit(ctx, "rate_limit:user123", 100, time.Minute)
if err != nil {
    if err.Error() == "unable to increment" {
        // Rate limit exceeded
    }
    // Handle other errors
}

// Circular counter with IncrWithMOD
// Useful for round-robin or rotating operations
// Returns values from 0 to (mod-1)
newVal, err := client.IncrWithMOD(ctx, "round_robin:workers", 3, time.Hour)
if err != nil {
    // Handle error
}
// newVal will cycle through 0, 1, 2, 0, 1, 2, ...
```

The counter operations provide atomic implementations for common rate limiting and circular counter patterns:

#### IncrWithLimit
- Atomically increments a counter with a maximum limit
- Sets expiry on first creation
- Returns error when limit is reached
- Useful for:
  - Rate limiting
  - Quota management
  - Concurrency control

#### IncrWithMOD
- Atomically increments a counter with modulo operation
- Optional expiry duration (0 means no expiry)
- Always returns values in range [0, mod-1]
- Useful for:
  - Round-robin distribution
  - Load balancing
  - Circular buffers

Both operations are implemented using Lua scripts to ensure atomicity and consistency.

## Best Practices

1. **Connection Management**:
   - Always close the client when done using `defer client.Close()`
   - Configure appropriate pool size based on your application needs

2. **Timeouts**:
   - Set appropriate timeouts for your use case
   - Default timeouts are:
     - DialTimeout: 500ms
     - ReadTimeout: 2000ms
     - WriteTimeout: 2000ms
     - IdleTimeout: 600s

3. **Cluster Mode**:
   - When using cluster mode, ensure all keys in multi-key operations belong to the same slot
   - Use `ServeReadsFromSlaves` or `ServeReadsFromMasterAndSlaves` for read scaling

4. **Error Handling**:
   - Always check for `redis.Nil` error when expecting key existence
   - Handle connection errors appropriately

## Performance Considerations

1. **Connection Pool**:
   - `PoolSize` is per node in cluster mode
   - Default pool size is 50 connections
   - Configure `MinIdleConn` to reduce connection establishment overhead

2. **Pipeline Operations**:
   - Use pipeline for batching multiple operations
   - Helps reduce network round trips

3. **Monitoring**:
   - Built-in New Relic integration for monitoring Redis operations
   - Track connection pool usage and operation latencies

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
