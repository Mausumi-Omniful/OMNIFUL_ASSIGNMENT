# Package monitoring

## Overview
The monitoring package provides a comprehensive solution for collecting and reporting application performance metrics, with a specific focus on queue monitoring and Prometheus integration. It offers tools to track various metrics including message production, consumption, and processing duration in distributed systems.

## Features
- **Prometheus Integration**: Built-in support for Prometheus metrics collection and exposure
- **Queue Monitoring**: Specialized monitoring for message queue operations
- **Custom Metrics**: Support for custom labels and metric prefixes
- **HTTP Endpoints**: Ready-to-use metrics endpoint for Prometheus scraping
- **Extensible Design**: Easy integration with multiple monitoring solutions

## Components

### 1. Queue Monitoring
The package provides detailed monitoring for queue operations with the following metrics:
- `queue_message_produced_total`: Tracks total messages produced to queues
- `queue_message_consumed_total`: Tracks total messages consumed from queues
- `queue_message_processing_duration`: Measures message processing duration

### 2. HTTP Endpoints
- `/metrics`: Exposes Prometheus-formatted metrics

## Installation
```bash
go get github.com/omniful/go_commons/monitoring
```

## Usage Examples

### 1. Basic Setup
```go
package main

import (
	"github.com/omniful/go_commons/monitoring"
	"github.com/omniful/go_commons/http"
)

func main() {
	// Create an HTTP server
	server := http.NewServer()
	
	// Register monitoring routes
	monitoring.RegisterMonitoringRoutes(server)
}
```

### 2. Queue Monitoring Setup
```go
package main

import (
	"github.com/omniful/go_commons/monitoring/queue"
	"time"
)

func main() {
	// Configure queue monitoring
	config := queue.MonitoringConfig{
		Prefix: "myapp",
		MessageProducedCustomLabels: []string{"service_name"},
		MessageConsumedCustomLabels: []string{"service_name"},
		MessageProcessingDurationCustomLabels: []string{"service_name"},
	}

	// Register a queue monitor
	err := queue.Register("myqueue", config)
	if err != nil {
		panic(err)
	}

	// Get the queue monitor
	monitor, err := queue.GetQueueMonitor("myqueue")
	if err != nil {
		panic(err)
	}

	// Record metrics
	attributes := queue.MonitoringAttributes{
		"service_name": "order_service",
	}

	// Record message received
	monitor.RecordMessageReceived("order_queue", attributes)

	// Record message consumed
	monitor.RecordMessageConsumed("order_queue", attributes)

	// Record processing duration
	monitor.RecordProcessingDuration("order_queue", attributes, 100*time.Millisecond)
}
```

## Configuration Options

### Queue Monitoring Config
```go
type MonitoringConfig struct {
	Prefix                                string   // Prefix for all metrics
	MessageProducedCustomLabels           []string // Custom labels for produced messages
	MessageConsumedCustomLabels           []string // Custom labels for consumed messages
	MessageProcessingDurationCustomLabels []string // Custom labels for processing duration
}
```

## Best Practices
1. **Consistent Labeling**: Use consistent label names across your application
2. **Meaningful Metrics**: Choose meaningful queue names and labels that provide valuable insights
3. **Error Handling**: Always handle errors returned by monitoring functions
4. **Resource Management**: Register monitors during application initialization

## Notes
- The package is designed for extensibility and integration with multiple monitoring solutions
- All metrics are automatically registered with Prometheus
- Queue names are automatically added as labels to all queue-related metrics
- Custom labels can be added to provide more context to your metrics

## Dependencies
- github.com/prometheus/client_golang
- github.com/gin-gonic/gin
- github.com/omniful/go_commons/http
- github.com/omniful/go_commons/log
- github.com/omniful/go_commons/util
