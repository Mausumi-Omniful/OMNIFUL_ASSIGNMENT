# Rate Limiter Package

## Overview
The ratelimiter package provides a flexible and robust rate limiting solution for Go applications, particularly designed for use with the Gin web framework. It offers external rate limiting capabilities that can be easily integrated into your web services to prevent abuse and ensure system stability.

## Features
- Configurable rate limits and expiry durations
- Customizable key extraction for identifying clients
- Middleware-based integration with Gin
- Thread-safe implementation
- Support for distributed rate limiting (when using a distributed storage backend)

## Components

### ExternalRateLimiter Interface
```go
type ExternalRateLimiter interface {
    INCRWithLimit(ctx context.Context, key string, limit int, expiry time.Duration) (int64, error)
}
```
This interface defines the contract for implementing rate limiting mechanisms. You can implement this interface with various backends (e.g., Redis, in-memory store).

### Options
```go
type Options struct {
    // Limit defines the default static rate limit
    Limit           int

    // Expiry defines the default time window duration
    Expiry          time.Duration

    // KeyExtractor extracts a unique identifier from the request for rate limiting
    KeyExtractor    func(*gin.Context) string

    // LimitExtractor dynamically determines the rate limit for each request
    // If not set, falls back to static Limit
    LimitExtractor  func(*gin.Context) int

    // ExpiryExtractor dynamically determines the time window duration for each request
    // If not set, falls back to static Expiry
    ExpiryExtractor func(*gin.Context) time.Duration
}
```

## Usage Examples

### Basic Usage with Gin Middleware
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/omniful/go_commons/ratelimiter/external_rate_limiter"
    "github.com/omniful/go_commons/redis"
    "github.com/omniful/go_commons/redis_cache"
    "github.com/omniful/go_commons/serializer"
    "github.com/omniful/go_commons/config"
    "time"
)

func main() {
    router := gin.Default()

    ctx := context.TODO()
    
    // Initialize your rate limiter implementation
    redisRateLimiter := cache.NewRedisCacheClient(
        redis.GetClient(), 
        cache.NewMsgpackSerializer(), 
        config.GetString(ctx, "service.name"),
    )
    
    // Configure middleware with options
    middleware := external_rate_limiter.ExternalRateLimitMiddleware(
        redisRateLimiter,
        external_rate_limiter.WithLimit(100),
        external_rate_limiter.WithExpiry(time.Minute),
        external_rate_limiter.WithKeyExtractor(func(c *gin.Context) string {
            return c.ClientIP() // Use client IP as the rate limit key
        }),
    )
    
    // Apply middleware to routes
    router.Use(middleware)
    
    // Your routes here
    router.GET("/api", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello, World!"})
    })
    
    router.Run(":8080")
}
```

### Dynamic Rate Limiting
```go
// Rate limit based on user ID from header
middleware := external_rate_limiter.ExternalRateLimitMiddleware(
    rateLimiter,
    external_rate_limiter.WithLimitExtractor(func(c *gin.Context) int {
        tier := c.GetHeader("X-Tier")
        switch tier {
        case "tier_3":
            return 1000
        case "tier_2":
            return 500
        case "tier_1":
            return 100
        default:
            return 10
        }
    }),
    external_rate_limiter.WithExpiryExtractor(func(c *gin.Context) time.Duration {
        tier := c.GetHeader("X-Tier")
        switch tier {
        case "tier_3":
            return time.Hour * 24
        case "tier_2":
            return time.Hour
        case "tier_1":
            return time.Minute * 30
        default:
            return time.Minute
        }
    }),
    external_rate_limiter.WithKeyExtractor(func(c *gin.Context) string {
        return c.GetHeader("X-Client-ID")
    }),
)
```

## Configuration Options

### Static Configuration
```go
// WithLimit sets a fixed rate limit
WithLimit(100) // Allow 100 requests per window

// WithExpiry sets the time window for the rate limit
WithExpiry(time.Minute) // Reset counter after 1 minute

// WithKeyExtractor defines how to identify clients
WithKeyExtractor(func(c *gin.Context) string {
    return c.ClientIP()
})
```

### Dynamic Configuration
```go
// WithLimitExtractor dynamically determines rate limit per request
WithLimitExtractor(func(c *gin.Context) int {
    // Example implementation
    tier := c.GetHeader("X-Tier")
    switch tier {
    case "tier_3":
        return 1000
    case "tier_2":
        return 500
    default:
        return 100
    }
})

// WithExpiryExtractor dynamically determines expiry duration per request
WithExpiryExtractor(func(c *gin.Context) time.Duration {
    // Example implementation
    tier := c.GetHeader("X-Tier")
    switch tier {
    case "tier_3":
        return time.Hour
    default:
        return time.Minute
    }
})
```

## Default Values
- Default Limit: 10 requests
- Default Expiry: 1 minute

## Error Handling
When a rate limit is exceeded, the middleware will:
1. Abort the request
2. Return HTTP 429 (Too Many Requests)
3. Include a JSON error message: `{"error": "Request Quota exceeded, Kindly slow down"}`

## Best Practices
1. Choose appropriate limits based on your application's requirements
2. Use dynamic rate limiting for different user tiers or service levels
3. Use distributed storage (e.g., Redis) for rate limiting in clustered environments
4. Consider implementing retry-after headers for better client experience
5. Monitor rate limiting metrics to adjust limits as needed

## Thread Safety
The implementation is thread-safe and suitable for concurrent use in production environments.
