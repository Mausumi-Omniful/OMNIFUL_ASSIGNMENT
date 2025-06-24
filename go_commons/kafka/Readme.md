# Package kafka

## Overview
The kafka package provides a robust interface for integrating with Apache Kafka in Go applications. It offers feature-rich producer and consumer implementations with support for synchronous and asynchronous message publishing, consumer groups, dead letter queues, and various authentication methods including SASL.

## Features
- **Flexible Producer**
  - Synchronous and asynchronous message publishing
  - Batch message support
  - Message compression (GZIP, Snappy)
  - Configurable acknowledgment levels
  - Automatic request ID propagation from context
  - FIFO ordering support through message keys
  
- **Consumer Groups**
  - Topic subscription with consumer groups
  - Message handler registration
  - Automatic rebalancing
  - Dead letter queue support for failed messages
  
- **Authentication Support**
  - SASL Plain authentication
  - IAM authentication for MSK
  - TLS support

- **Additional Features**
  - Message headers support with automatic request ID injection
  - Request ID tracking and propagation
  - Interceptor support (e.g., NewRelic)
  - Configurable retry intervals
  - Comprehensive error handling

## Configuration
The package uses a flexible configuration system with functional options:

```go
type Config struct {
    Brokers          []string
    ConsumerGroupID  string
    TransactionName  string
    RetryInterval    time.Duration
    ClientID         string
    Version          string
    Region           string
    DeadLetterQueue  *DeadLetterQueueConfig
}

type DeadLetterQueueConfig struct {
    Queue     string
    Account   string
    Endpoint  string
    Prefix    string
    ShouldLog bool
    Region    string
}
```

## Usage Examples

### Producer Example with Message Keys and Headers
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/kafka"
    "github.com/omniful/go_commons/pubsub"
)

func main() {
    // Initialize producer with configuration
    producer := kafka.NewProducer(
        kafka.WithBrokers([]string{"localhost:9092"}),
        kafka.WithClientID("my-producer"),
        kafka.WithKafkaVersion("2.8.1"),
    )
    defer producer.Close()

    // Create message with key for FIFO ordering
    msg := &pubsub.Message{
        Topic: "my-topic",
        // Key is crucial for maintaining FIFO ordering
        // Messages with the same key will be delivered to the same partition in order
        Key: "customer-123",  
        Value: []byte("Hello Kafka!"),
        Headers: map[string]string{
            "custom-header": "value",
            // Note: HeaderXOmnifulRequestID will be automatically added
            // from context if present
        },
    }

    // Context with request ID
    ctx := context.WithValue(context.Background(), "request_id", "req-123")
    
    // Synchronous publish - HeaderXOmnifulRequestID will be automatically added
    err := producer.Publish(ctx, msg)
    if err != nil {
        panic(err)
    }

    // Batch publish with consistent keys for ordering
    messages := []*pubsub.Message{
        {
            Topic: "my-topic",
            Key: "customer-123",  // Same key maintains ordering
            Value: []byte("Message 1"),
        },
        {
            Topic: "my-topic",
            Key: "customer-123",  // Same key maintains ordering
            Value: []byte("Message 2"),
        },
    }
    err = producer.PublishBatch(ctx, messages)
    if err != nil {
        panic(err)
    }
}
```

### Consumer Example with Interceptor
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/kafka"
    "github.com/omniful/go_commons/pubsub"
    "github.com/omniful/go_commons/pubsub/interceptor"
)

// Implement message handler
type MessageHandler struct{}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
    // Process message
    return nil
}

func main() {
    // Initialize consumer with configuration
    consumer := kafka.NewConsumer(
        kafka.WithBrokers([]string{"localhost:9092"}),
        kafka.WithConsumerGroup("my-consumer-group"),
        kafka.WithClientID("my-consumer"),
        kafka.WithKafkaVersion("2.8.1"),
        kafka.WithRetryInterval(time.Second),
        kafka.WithDeadLetterConfig(&kafka.DeadLetterQueueConfig{
            Queue:     "dlq-queue",
            Account:   "aws-account",
            Region:    "us-east-1",
            ShouldLog: true,
        }),
    )
    defer consumer.Close()

    // Set NewRelic interceptor for monitoring
    consumer.SetInterceptor(interceptor.NewRelicInterceptor())

    // Register message handler for topic
    handler := &MessageHandler{}
    consumer.RegisterHandler("my-topic", handler)

    // Start consuming messages
    ctx := context.Background()
    consumer.Subscribe(ctx)
}
```

### Authentication Examples

#### SASL Plain Authentication
```go
producer := kafka.NewProducer(
    kafka.WithBrokers([]string{"localhost:9092"}),
    kafka.WithSASLPlainAuthentication("username", "password"),
)
```

#### IAM Authentication (for MSK)
```go
producer := kafka.NewProducer(
    kafka.WithBrokers([]string{"localhost:9092"}),
    kafka.WithRegion("us-east-1"),
    kafka.WithIAMAuthentication(true),
)
```

## Best Practices
1. Always close producers and consumers when done
2. Use appropriate retry intervals for your use case
3. Implement proper error handling in message handlers
4. Consider using dead letter queues for failed messages
5. Use batch publishing for better performance when sending multiple messages
6. Configure appropriate acknowledgment levels based on your reliability requirements
7. Use message keys consistently for maintaining FIFO ordering:
   - Messages with the same key are guaranteed to be processed in order
   - Different keys allow for parallel processing
   - Empty keys result in round-robin partition distribution
8. Be aware that request IDs from context are automatically propagated to message headers

## Message Ordering and Keys
- Messages with the same key are guaranteed to be delivered to the same partition in order
- The Key field in `pubsub.Message` determines the partition assignment
- Use cases for message keys:
  - Customer ID: Ensure all messages for a customer are processed in order
  - Order ID: Maintain sequence of order-related events
  - Account ID: Keep account transactions in order
- Empty keys result in round-robin distribution across partitions (no ordering guarantee)

## Request ID Propagation
- The package automatically injects `HeaderXOmnifulRequestID` into message headers
- The request ID is extracted from the context if present
- This enables end-to-end request tracing across services
- Useful for:
  - Distributed tracing
  - Debugging
  - Request correlation
  - Monitoring and observability

## Interceptors
The package provides a powerful interceptor mechanism that allows you to add cross-cutting concerns to your message processing pipeline. Interceptors can be used to add functionality before or after producing/consuming messages.

### Built-in NewRelic Interceptor
The package includes a built-in NewRelic interceptor that provides:
- Automatic transaction creation and naming
- Request ID propagation to NewRelic transactions
- Error tracking and reporting to NewRelic
- Panic recovery with proper error reporting
- Performance monitoring of message processing

### Custom Interceptors
You can create custom interceptors by implementing the `Interceptor` interface:

```go
type Interceptor func(ctx context.Context, msg *pubsub.Message, handler Handler) error
```

Example of a custom logging interceptor:
```go
func LoggingInterceptor() interceptor.Interceptor {
    return func(ctx context.Context, msg *pubsub.Message, handler interceptor.Handler) error {
        log.Infof("Processing message from topic: %s", msg.Topic)
        err := handler(ctx, msg)
        if err != nil {
            log.Errorf("Error processing message: %v", err)
        }
        return err
    }
}
```

### Chaining Interceptors
Multiple interceptors can be chained together using the `ChainInterceptors` function:

```go
consumer.SetInterceptor(interceptor.ChainInterceptors(
    LoggingInterceptor(),
    interceptor.NewRelicInterceptor(),
    // Add more interceptors...
))
```

### Interceptor Best Practices
1. Use interceptors for cross-cutting concerns like:
   - Monitoring and metrics
   - Logging and tracing
   - Error handling
   - Authentication/Authorization
2. Keep interceptors focused on a single responsibility
3. Order interceptors properly when chaining (e.g., logging before monitoring)
4. Handle errors appropriately in interceptors
5. Be mindful of performance impact when adding multiple interceptors

## Notes
- The package uses the [Sarama](https://github.com/IBM/sarama) library internally
- Supports Kafka versions compatible with Sarama
- Includes monitoring integration capabilities
- Thread-safe for concurrent usage
