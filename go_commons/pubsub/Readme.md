# PubSub Package

This package provides a generic publish-subscribe (pub/sub) messaging interface for Go applications. It offers a flexible and extensible way to implement message-based communication patterns in your applications.

## Overview

The package consists of three main components:

1. **Publisher Interface**: For publishing messages to topics
2. **Subscriber Interface**: For subscribing to topics and handling messages
3. **Message Structure**: For encapsulating message data and metadata

## Components

### Publisher Interface

```go
type Publisher interface {
    Publish(ctx context.Context, msg *Message) error
    PublishAsync(ctx context.Context, msg *Message) error
    PublishBatch(ctx context.Context, msgs []*Message) error
    Close()
}
```

The Publisher interface provides methods for:
- Synchronous message publishing
- Asynchronous message publishing
- Batch message publishing
- Graceful shutdown

### Subscriber Interface

```go
type Subscriber interface {
    Subscribe(ctx context.Context) error
    RegisterHandler(topic string, handler IPubSubMessageHandler)
    Close()
}
```

The Subscriber interface allows:
- Subscribing to topics
- Registering message handlers
- Graceful shutdown

### Message Structure

```go
type Message struct {
    Topic     string
    Value     []byte
    Key       string
    Headers   map[string]string
    Timestamp time.Time
}
```

The Message struct includes:
- Topic name
- Message payload (Value)
- Message key for routing/partitioning
- Custom headers
- Timestamp information

## Usage Examples

### Creating and Publishing a Message

```go
package main

import (
    "context"
    "time"
    "github.com/omniful/go_commons/pubsub"
)

func main() {
    // Create a new message
    msg := &pubsub.Message{
        Topic:     "orders",
        Key:       "order-123",
        Headers:   map[string]string{"version": "1.0"},
        Timestamp: time.Now(),
    }

    // Create payload
    type Order struct {
        ID     string `json:"id"`
        Amount float64 `json:"amount"`
    }
    
    order := Order{
        ID:     "order-123",
        Amount: 99.99,
    }

    // Convert payload to bytes
    payload, err := pubsub.NewEventInBytes(order)
    if err != nil {
        panic(err)
    }
    msg.Value = payload

    // Publish message (assuming you have a publisher implementation)
    publisher.Publish(context.Background(), msg)
}
```

### Implementing a Message Handler

```go
package main

import (
    "context"
    "github.com/omniful/go_commons/pubsub"
)

type OrderHandler struct{}

func (h *OrderHandler) Process(ctx context.Context, msg *pubsub.Message) error {
    // Create a struct to unmarshal the message into
    var order struct {
        ID     string  `json:"id"`
        Amount float64 `json:"amount"`
    }

    // Unmarshal the message
    if err := msg.To(&order, ""); err != nil {
        return err
    }

    // Process the order
    // ... your business logic here ...

    return nil
}

func main() {
    // Register the handler (assuming you have a subscriber implementation)
    handler := &OrderHandler{}
    subscriber.RegisterHandler("orders", handler)
    
    // Start subscribing
    ctx := context.Background()
    if err := subscriber.Subscribe(ctx); err != nil {
        panic(err)
    }
}
```

## Best Practices

1. **Error Handling**: Always handle errors returned by Publish and Subscribe operations.
2. **Context Usage**: Use context for timeout and cancellation control.
3. **Graceful Shutdown**: Always call Close() on publishers and subscribers when shutting down.
4. **Message Validation**: Validate messages before publishing and after receiving.
5. **Handler Implementation**: Keep message handlers lightweight and non-blocking.

## Implementation Notes

- The package provides interfaces only. You need to implement these interfaces for your specific message broker (e.g., Kafka, RabbitMQ, Redis).
- Message handlers should be thread-safe as they might be called concurrently.
- The `To()` method on Message provides type-safe deserialization of message payloads.
- Async publishing (`PublishAsync`) should be used for high-throughput scenarios.

## Thread Safety

The interfaces are designed to be thread-safe, but specific implementations must ensure thread safety for:
- Message publishing
- Handler registration
- Subscription management
