# Package interservice-client

## Overview
The interservice-client package provides a robust and flexible solution for inter-service communication in Go applications. It offers a high-level abstraction for making HTTP requests between microservices, with built-in support for error handling, request/response processing, and configurable transport mechanisms.

## Features
- Simple and intuitive API for HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Built-in response parsing and error handling
- Configurable HTTP transport settings
- Support for request validation
- NewRelic instrumentation
- Context-aware request handling
- Customizable timeout and retry mechanisms
- Structured error responses

## Installation
```bash
go get github.com/omniful/go_commons/interservice-client
```

## Configuration
The client can be configured using the `Config` struct:

```go
type Config struct {
    ServiceName string         // Name of the service
    BaseURL     string         // Base URL for the service
    Timeout     time.Duration  // Request timeout duration
    Transport   *http.Transport // Custom transport configuration (optional)
}
```

### Basic Setup
```go
package main

import (
    "time"
    "github.com/omniful/go_commons/interservice-client"
)

func main() {
    config := interservice_client.Config{
        ServiceName: "user-service",
        BaseURL:     "http://user-service.internal",
        Timeout:     5 * time.Second,
    }
    
    client, err := interservice_client.NewClientWithConfig(config)
    if err != nil {
        panic(err)
    }
    // Client is ready to use
}
```

## Usage Examples

### Making a GET Request
```go
func GetUserDetails(ctx context.Context, userID string) (*UserDetails, error) {
    request := &http.Request{
        URL: fmt.Sprintf("/users/%s", userID),
    }
    
    var userData UserDetails
    response, err := client.Get(request, &userData)
    if err != nil {
        return nil, err
    }
    
    return &userData, nil
}
```

### Making a POST Request with Data
```go
func CreateUser(ctx context.Context, user *User) error {
    request := &http.Request{
        URL:  "/users",
        Body: user,
    }
    
    var result CreateUserResponse
    response, err := client.Post(request, &result)
    if err != nil {
        return err
    }
    
    return nil
}
```

### Using Validation
```go
func CreateUserWithValidation(ctx context.Context, user *User) error {
    request := &http.Request{
        URL:  "/users",
        Body: user,
    }
    
    validator := validatorPkg.New()
    var result CreateUserResponse
    response, err := client.ExecuteWithValidation(
        ctx,
        http.APIPost,
        request,
        &result,
        validator,
    )
    if err != nil {
        return err
    }
    
    return nil
}
```

## Error Handling
The package provides structured error handling with the `Error` type:

```go
type Error struct {
    Message    string            // Error message
    Errors     map[string]string // Detailed error information
    StatusCode http.StatusCode   // HTTP status code
    Data       interface{}       // Additional error data
}
```

### Handling Errors in Gin Handlers
```go
func UserHandler(c *gin.Context) {
    response, err := client.Get(request, &data)
    if err != nil {
        interservice_client.NewErrorResponseByInterServiceError(c, err)
        return
    }
    // Process successful response
}
```

## Response Structure
All responses follow a consistent structure:

```go
type InterSvcResponse struct {
    IsSuccess  bool                   // Success status
    StatusCode int                    // HTTP status code
    Data       json.RawMessage        // Response data
    Meta       map[string]interface{} // Metadata
    Error      json.RawMessage        // Error information
}
```

## Best Practices
1. Always use context for request timeouts and cancellation
2. Handle errors appropriately using the provided error types
3. Use validation when dealing with structured data
4. Configure appropriate timeouts for your use case
5. Monitor request metrics using the built-in NewRelic instrumentation

## Notes
- The client automatically handles common scenarios like context cancellation and timeouts
- Default transport settings are optimized for microservice communication
- The client supports structured error responses for better error handling
- Built-in support for request tracing and monitoring
