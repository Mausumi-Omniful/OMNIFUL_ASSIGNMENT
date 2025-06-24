# SQS Package Documentation

The SQS package provides a robust abstraction for interacting with AWS SQS. It simplifies both the production and consumption of messages in distributed systems by offering clear APIs for sending data, batching messages, and concurrently processing messages with a worker pool.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [File Overview](#file-overview)
- [Usage](#usage)
  - [Publisher Usage](#publisher-usage)
  - [Consumer Usage](#consumer-usage)
- [Working Example](#working-example)
- [Configuration](#configuration)
- [Environment Variables for Local Testing](#environment-variables-for-local-testing)
- [Questions & Clarifications](#questions--clarifications)
- [Contributing](#contributing)
- [License](#license)

## Overview
The sqs package is designed to enable developers to easily integrate AWS SQS messaging into their Go applications. It encapsulates functionalities for:

- Publishing messages to an SQS queue
- Batching messages to improve efficiency
- Setting up consumers to continuously poll and process messages
- Managing worker pools for concurrent message processing

By abstracting the low-level AWS SQS interactions, the package helps streamline the messaging workflows in distributed systems.

## Features
- **Publisher API:** Simplified methods to publish single or batched messages.
- **Batch Processing:** Utilities to group messages and send them in bulk, automatically handling SQS size limits.
- **Consumer API:** Tools for consuming messages with support for concurrent processing.
- **Worker Pool:** Management of workers to optimize processing throughput.
- **Handler Wrapping:** Built-in middleware for authentication, logging, and NewRelic metrics.
- **Automatic Compression:** Smart compression handling for large messages.

## File Overview

### producer.go
Contains the implementation of the `Publisher` struct for sending messages to an SQS queue. It handles both individual and batched message publishing, ensuring optimized interaction with the SQS API. Includes automatic compression for messages exceeding SQS size limits.

### batch.go
Provides structures and functions for handling SQS message size limits (256KB). It automatically breaks down large messages into acceptable sizes while maintaining message integrity.

### queue.go
Defines the core functionalities to initialize and manage SQS queues, including configuration settings and queue-specific operations.

### wrapper_handler.go
Implements a wrapper around user-defined message handlers to provide built-in middleware functionality including:
- Authentication handling
- Logging
- NewRelic metrics tracking
This ensures consistent handling of these cross-cutting concerns without developer intervention.

### consumer_worker.go
Implements the worker logic that continuously fetches messages from SQS and processes them using a designated handler. It facilitates concurrent message processing.

### message.go
Defines the Message struct that standardizes the message format used throughout the package. This abstraction simplifies message creation and retrieval.

### consumer.go
Sets up the consumer that continuously polls the specified SQS queue for new messages. It utilizes the wrapped handler to process each message as it arrives.

### pool.go
Manages a pool of consumer workers, balancing the load and ensuring that multiple messages can be processed in parallel.

### worker.go
Contains supporting functions and routines for the worker mechanism, including error handling and message acknowledgment processes.

### handler.go
Defines the interfaces and contracts for message handlers. Developers can implement these interfaces to create custom handlers that integrate into the consumer workflow.

### sqs.go
Acts as the main entry point for the package. It ties together publisher and consumer functionalities and may include initialization routines and shared configurations.

## Usage

### Message Compression

The package provides flexible compression options that can be configured at both queue and message levels:

1. **Queue-Level Compression:**
```go
queue := &sqs.Queue{
    // ... other configuration ...
    Compressor: compression.GetCompressionParser(compression.GZIP),
}
```

2. **Message-Level Compression:**
```go
msg := &sqs.Message{
    Value: []byte("Large payload"),
    Compression: compression.GZIP,
}
```

3. **Automatic Compression:**
The package automatically handles compression for messages exceeding SQS size limits (250KB):
```go
publisher := sqs.NewPublisher(queue)
msg := &sqs.Message{
    Value: []byte("Very large payload"), // If > 250KB, automatically compressed
}
```

### Publisher Usage

To send messages to SQS:
1. Create a new publisher instance using the `NewPublisher` function.
2. Construct a message using the Message struct.
3. Optionally, batch messages if sending multiple messages simultaneously.
4. Use the `Publish` or `BatchPublish` methods to deliver the message(s) to the queue.

**Example:**
```go
queue := // Initialize Queue object
publisher := sqs.NewPublisher(queue)
msg := &sqs.Message{
    Value: []byte("Your message payload"),
    // Set other message properties as needed
}

ctx := context.Background()
if err := publisher.Publish(ctx, msg); err != nil {
    log.Fatalf("Failed to publish message: %v", err)
}
```

### Batch Processing

The package automatically handles SQS message size limits (256KB) when batching messages:

```go
messages := []*sqs.Message{
    {Value: []byte("Message 1")},
    {Value: []byte("Message 2")},
    // ... more messages
}

ctx := context.Background()
if err := publisher.BatchPublish(ctx, messages); err != nil {
    log.Fatalf("Failed to publish batch: %v", err)
}
```

The batch processing will:
- Automatically split messages if they exceed SQS size limits
- Handle compression for large messages
- Maintain message integrity throughout the process

### Consumer Usage

To set up a consumer:
1. Define a message handler that implements the `ISqsMessageHandler` interface
2. Create a consumer instance using the `NewConsumer` function with appropriate configuration
3. Start the consumer with a context to begin processing incoming messages

**Example:**
```go
type MyHandler struct{}

func (h *MyHandler) Handle(msg *sqs.Message) error {
    fmt.Println("Processing message:", string(msg.Value))
    return nil
}

queue := // Initialize Queue object
handler := &MyHandler{}

consumer, err := sqs.NewConsumer(
    queue,           // Queue configuration
    1,              // Number of workers
    1,              // Concurrency per worker
    handler,        // Message handler
    10,             // Max messages count
    30,             // Visibility timeout
    false,          // Is async
    false,          // Send batch message
)
if err != nil {
    log.Fatalf("Failed to create consumer: %v", err)
}

ctx := context.Background()
consumer.Start(ctx)

// To gracefully shutdown the consumer
// consumer.Close()
```

## Working Example

Below is a complete working example that demonstrates both publishing and consuming messages using the sqs package:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/omniful/go_commons/sqs"
)

type ExampleHandler struct{}

func (h *ExampleHandler) Handle(msg *sqs.Message) error {
    fmt.Println("Processing message:", string(msg.Value))
    return nil
}

func main() {
    // Initialize AWS session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    })
    if err != nil {
        log.Fatalf("Failed to create AWS session: %v", err)
    }

    queueURL := "https://sqs.us-east-1.amazonaws.com/your-account/your-queue"
    
    // Initialize Queue
    queue := // Initialize Queue object with sess and queueURL

    // Set up publisher
    publisher := sqs.NewPublisher(queue)
    message := &sqs.Message{
        Value: []byte("Hello SQS!"),
    }

    ctx := context.Background()
    if err := publisher.Publish(ctx, message); err != nil {
        log.Fatalf("Failed to publish message: %v", err)
    }

    // Set up consumer
    handler := &ExampleHandler{}
    consumer, err := sqs.NewConsumer(
        queue,
        1,       // Number of workers
        1,       // Concurrency per worker
        handler,
        10,      // Max messages count
        30,      // Visibility timeout
        false,   // Is async
        false,   // Send batch message
    )
    if err != nil {
        log.Fatalf("Failed to create consumer: %v", err)
    }

    consumer.Start(ctx)

    // Let the consumer run for a while
    time.Sleep(10 * time.Second)
    consumer.Close()
}
```

## Configuration

The package supports various configurations including:

- AWS credentials and region via the AWS SDK's session configuration
- SQS specific settings such as queue URL, maximum batch size, and polling intervals
- Compression settings:
  - Queue-level default compression
  - Message-specific compression
  - Automatic compression for large messages
- Consumer options:
  - Number of workers
  - Concurrency per worker
  - Max messages count (up to 10)
  - Visibility timeout
  - Async processing
  - Batch message processing

## Environment Variables for Local Testing

The package supports the following environment variables that are particularly useful during local development and testing:

### LOCAL_SQS_ENDPOINT
- **Purpose**: Sets a custom endpoint for SQS connections
- **Usage**: Set this variable when using LocalStack or other AWS service emulators for local testing
- **Example**: `export LOCAL_SQS_ENDPOINT=http://localhost:4566`

### AWS_DEBUG_LOG
- **Purpose**: Enables detailed AWS SQS operation logging
- **Usage**: Set to "true" to enable debug logging of AWS SQS operations
- **Example**: `export AWS_DEBUG_LOG=true`

## Questions & Clarifications

## Contributing

Contributions are welcome! Please follow the Go community guidelines when submitting issues or pull requests. Ensure that new code is well tested and documented.
