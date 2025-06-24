# Package httpclient

## Overview
The `httpclient` package provides a powerful, feature-rich HTTP client implementation designed for robust and flexible HTTP communications in Go applications. Built on top of the `resty` client, it offers extensive configuration options, middleware support, and enterprise-ready features out of the box.

## Key Features
- **Flexible Configuration**: Extensive options for customizing client behavior
- **Retry Mechanisms**: Built-in support for configurable retry strategies
- **Rate Limiting**: Pluggable rate limiting capabilities
- **Authentication**: Support for both client-level and request-level authentication
- **Middleware Support**: Pre and post-request hooks via callback system
- **Panic Recovery**: Built-in panic handling with customizable recovery strategies
- **Logging**: Configurable request/response logging
- **Context Support**: Full integration with Go's context for cancellation and deadlines
- **Request ID Tracking**: Automatic request ID propagation

## Installation
```bash
go get github.com/omniful/go_commons
```

## Basic Usage
```go
package main

import (
	"context"
	"github.com/omniful/go_commons/httpclient"
	"github.com/omniful/go_commons/httpclient/request"
)

func main() {
	// Create a new client with base URL
	client := httpclient.New("https://api.example.com")

	// Create a request
	req, _ := request.NewBuilder().
		SetPath("/users").
		SetMethod("GET").
		Build()

	// Send request
	resp, err := client.Get(context.Background(), req)
	if err != nil {
		panic(err)
	}

	// Process response
	var users []User
	if err := resp.Unmarshal(&users); err != nil {
		panic(err)
	}
}
```

## Advanced Configuration
```go
client := httpclient.New(
	"https://api.example.com",
	// Set client-level authentication
	httpclient.WithClientAuth(httpclient.BearerAuth("your-token")),
	
	// Configure retries
	httpclient.WithRetry(NewDefaultRetry(3)),
	
	// Set rate limiting
	httpclient.WithRateLimiter(NewTokenBucketRateLimiter(10, time.Second)),
	
	// Configure logging
	httpclient.WithLogConfig(httpclient.LogConfig{
		LogRequest:  true,
		LogResponse: true,
		LogLevel:    "info",
	}),
	
	// Set request timeout
	httpclient.WithDeadline(time.Second * 30),
	
	// Add request/response hooks
	httpclient.WithBeforeSendCallback(func(ctx *Context, req Request) error {
		// Pre-request logic
		return nil
	}),
	httpclient.WithAfterSendCallback(func(ctx *Context, req Request) error {
		// Post-request logic
		return nil
	}),
)
```

## Authentication
The package supports multiple authentication methods:

```go
// Bearer Token Authentication
client := httpclient.New(
	"https://api.example.com",
	httpclient.WithClientAuth(httpclient.BearerAuth("your-token")),
)

// Basic Authentication
client := httpclient.New(
	"https://api.example.com",
	httpclient.WithClientAuth(httpclient.BasicAuth("username", "password")),
)

// Custom Authentication Provider
type CustomAuthProvider struct{}

func (p *CustomAuthProvider) Apply(ctx context.Context, req Request) (Request, error) {
	// Add custom auth logic
	return req, nil
}

client := httpclient.New(
	"https://api.example.com",
	httpclient.WithRequestAuthProvider(&CustomAuthProvider{}),
)
```

## Retry Strategies
The package includes built-in retry support with customizable strategies:

```go
// Simple retry with max attempts
client := httpclient.New(
	"https://api.example.com",
	httpclient.WithMaxRetries(3),
)

// Custom retry strategy
type CustomRetry struct{}

func (r *CustomRetry) ShouldRetry(ctx context.Context, req Request, resp Response) bool {
	return resp.StatusCode() >= 500
}

func (r *CustomRetry) NextAttemptIn(ctx context.Context, req Request, resp Response) time.Duration {
	return time.Second * 2
}

client := httpclient.New(
	"https://api.example.com",
	httpclient.WithRetry(&CustomRetry{}),
)
```

## Error Handling
The package provides comprehensive error handling:

```go
client := httpclient.New(
	"https://api.example.com",
	httpclient.WithOnErrorCallback(func(ctx *Context, req Request, resp Response, err error) {
		// Custom error handling logic
	}),
	httpclient.WithPanicHandler(func(err error) {
		// Custom panic recovery logic
	}),
)
```

## Best Practices
1. Always use context for request cancellation and timeouts
2. Configure appropriate retry strategies for your use case
3. Implement proper error handling
4. Use request/response logging in development
5. Configure rate limiting for API compliance
6. Reuse client instances instead of creating new ones

## Thread Safety
The client is safe for concurrent use by multiple goroutines.

## Notes
- The client is built on top of the `resty` HTTP client
- All methods are context-aware and support cancellation
- Request IDs are automatically propagated unless disabled
- Default configuration includes reasonable timeouts and retry settings
