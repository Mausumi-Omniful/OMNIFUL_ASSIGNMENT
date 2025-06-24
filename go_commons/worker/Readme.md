# Worker Package

## Overview
The worker package provides a robust framework for managing and executing background tasks and concurrent job processing in Go applications. It implements a flexible, registry-based system for initializing and managing multiple types of listeners (HTTP, Kafka, SQS) with built-in graceful shutdown handling.

## Usage Guide

### 1. Configuration Setup (sample config.yaml)
```yaml
# Global AWS Configuration
aws:
  region: "eu-central-1"
  account: "348639420040"
  shouldLog: false
  sqs:
    prefix: "staging"  # Added prefix for SQS queue names
    endpoint: "https://sqs.eu-central-1.amazonaws.com/348639420040"

# Global Kafka Configuration
kafka:
  brokers:
    - "localhost:9092"
  clientId: oms-service
  version: 2.0.0
  iamAuthentication: false
  region: "eu-central-1"

# SQS Worker Configuration
workers:
  catalogueServiceCategoryQueue:
    name: "staging-product-service-category"
    workerCount: 1
    concurrencyPerWorker: 2
	workerGroup: "salesChannel"
   /// Add similar configs for other workers

# Kafka Consumer Configuration
consumers:
  updateOrder:
    name: "scs-order-update-consumer"
    workerGroup: "salesChannel"
    topic: "omniful.oms-service.order-events"
    groupId: "omniful.wms-service.order-events.local-scs-service.cg"
    enabled: true
	/// Add similar configs for other workers
```

### 2. Implementation in Your Service

```go
package main

import (
    "context"
    "github.com/omniful/go_commons/worker"
    "github.com/omniful/go_commons/worker/configs"
    "github.com/omniful/go_commons/worker/registry"
    "github.com/omniful/go_commons/pubsub"
    "github.com/omniful/go_commons/sqs"
    commonhttp "github.com/omniful/go_commons/http"
)

func main() {
    ctx := context.Background()
    
    // Initialize registry
    reg := registry.NewRegistry()
    
    // Register all available listeners
    reg.RegisterKafkaListenerConfig(
        ctx,
        "updateOrder",
        func(ctx context.Context, config configs.KafkaConsumerConfig) pubsub.IPubSubMessageHandler {
            return &KafkaMessageHandler{}
        },
    )

    reg.RegisterSQSListenerConfig(
        ctx,
        "catalogueServiceCategoryQueue",
        func(ctx context.Context, config configs.SqsQueueConfig) sqs.ISqsMessageHandler {
            return &SQSMessageHandler{}
        },
    )

    httpServer := commonhttp.NewServer()
    reg.RegisterHTTPListenerConfig(httpServer, "service-name")
    
    // Create server
    server := worker.NewServerFromRegistry(reg)
    
    // Option 1: Run all registered workers
    server.Run(ctx)
    
    // Option 2: Run specific worker groups
    serverConfig := configs.ServerConfig{
        IncludeGroupsArg: "salesChannel",  // Only run workers in salesChannel group
    }
    server.RunFromConfig(ctx, serverConfig)
    
    // Option 3: Run specific workers by name
    serverConfig = configs.ServerConfig{
        ListenerNamesArg: "scs-order-update-consumer,another-listener",  // Run specific workers
    }
    server.RunFromConfig(ctx, serverConfig)
    
    // Option 4: Run all workers except specific groups
    serverConfig = configs.ServerConfig{
        ExcludeGroupsArg: "salesChannel",  // Run all workers except those in salesChannel group
    }
    server.RunFromConfig(ctx, serverConfig)
}
```

### 3. Controlling Workers with ServerConfig

The ServerConfig provides flexible control over which workers to run. You can configure it either programmatically or through command-line arguments.

#### Command-Line Configuration
```go
package main

import (
    "context"
    "flag"
    "strings"
    "github.com/omniful/go_commons/worker"
    "github.com/omniful/go_commons/worker/configs"
    "github.com/omniful/go_commons/worker/registry"
    "github.com/omniful/go_commons/util"
)

func main() {
    var serverConfig configs.ServerConfig

    // Define command-line flags
    flag.StringVar(
        &serverConfig.IncludeGroupsArg,
        "includeGroups",
        "",
        "Comma-separated list of worker groups to include",
    )

    flag.StringVar(
        &serverConfig.ExcludeGroupsArg,
        "excludeGroups",
        "",
        "Comma-separated list of worker groups to exclude",
    )

    flag.StringVar(
        &serverConfig.ListenerNamesArg,
        "listenerNames",
        "",
        "Comma-separated list of specific listeners to run",
    )

    flag.Parse()

    // Initialize registry and server
    ctx := context.Background()
    reg := registry.NewRegistry()
    // ... register your listeners ...
    
    server := worker.NewServerFromRegistry(reg)
    server.RunFromConfig(ctx, serverConfig)
}
```

#### Controlling Workers Through Command-Line Arguments

1. **Run All Workers**
```bash
# Run all workers (no flags)
./your-service

# Or explicitly with empty flags
./your-service --includeGroups="" --excludeGroups="" --listenerNames=""
```

2. **Run Specific Worker Groups**
```bash
# Run only salesChannel and inventory workers
./your-service --includeGroups=salesChannel,inventory

# Run multiple groups
./your-service --includeGroups=salesChannel,inventory,catalogueSync
```

3. **Run Specific Workers**
```bash
# Run one or more specific workers
./your-service --listenerNames=scs-order-update-consumer,another-listener
```

4. **Exclude Specific Groups**
```bash
# Run all workers except salesChannel
./your-service --excludeGroups=salesChannel

# Exclude multiple groups
./your-service --excludeGroups=salesChannel,inventory
```

This approach allows you to:
- Control workers through command-line arguments
- Run different worker combinations based on needs
- Easily manage worker groups
- Enable/disable specific workers as needed

## Best Practices
1. Always define global configurations (`aws.*`, `kafka.*`) at the root level
2. Use meaningful names for workers and consumers that reflect their purpose
3. Configure appropriate worker counts and concurrency based on your needs
4. Enable/disable consumers using the `enabled` flag
5. Group related consumers using `workerGroup`
6. Use descriptive topic names and group IDs that follow your naming convention
7. Ensure listener names are unique across all types (HTTP, Kafka, SQS)

## Notes
- Ensure all required configurations are present in your `config.yaml`
- Test your worker configurations in a staging environment first
- Monitor worker performance and adjust concurrency settings as needed
- Implement proper error handling in your message processing functions
- Use context for proper cancellation and timeout handling

## Environment Variables for Local SQS Testing

### LOCAL_SQS_ENDPOINT
- **Purpose**: Sets a custom endpoint for SQS connections
- **Usage**: Set this variable when using LocalStack or other AWS service emulators for local testing
- **Example**: `export LOCAL_SQS_ENDPOINT=http://localhost:4566`

### AWS_DEBUG_LOG
- **Purpose**: Enables detailed AWS SQS operation logging
- **Usage**: Set to "true" to enable debug logging of AWS SQS operations
- **Example**: `export AWS_DEBUG_LOG=true`

## Contributing
Contributions are welcome! Please ensure:
- Code follows the project's style guidelines
- Tests are included for new functionality
- Documentation is updated for any changes