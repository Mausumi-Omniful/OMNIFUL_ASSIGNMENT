# Package pagination

## Overview
The pagination package provides a simple and efficient way to implement pagination in Go applications using Gin web framework and GORM. It includes middleware for handling pagination parameters from HTTP requests and utilities for applying pagination to database queries.

## Features
- HTTP middleware for handling pagination parameters
- GORM pagination utilities
- Configurable page size with sensible defaults
- Built-in parameter validation and sanitization

## Components

### Constants
```go
const (
    PerPage = "per_page"  // Query parameter for items per page
    Page    = "page"      // Query parameter for page number
    Limit   = "limit"     // Context key for limit value
)
```

### Middleware
The package provides a Gin middleware that automatically extracts and validates pagination parameters from HTTP requests:

```go
func main() {
    router := gin.New()
    
    // Apply pagination middleware to routes
    router.Use(pagination.Middleware())
    
    // Your routes here...
}
```

The middleware:
- Extracts `page` and `per_page` from query parameters
- Validates and sanitizes the values
- Sets default values if parameters are missing or invalid
- Stores the values in the Gin context for later use

### Database Utilities
The package includes GORM utilities for applying pagination to database queries:

```go
func GetUsers(c *gin.Context, db *gorm.DB) ([]User, error) {
    var users []User
    
    result := db.Scopes(pagination.Paginate(c)).Find(&users)
    return users, result.Error
}
```

## Configuration

### Page Size Limits
- Default page size: 20 items
- Maximum page size: 100 items
- Minimum page: 1

## Complete Example

Here's a complete example showing how to use the pagination package in a Gin API:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/omniful/go_commons/pagination"
    "gorm.io/gorm"
)

type User struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

func main() {
    // Initialize your GORM DB
    var db *gorm.DB
    
    // Create Gin router
    router := gin.New()
    
    // Apply pagination middleware
    router.Use(pagination.Middleware())
    
    // Example route with pagination
    router.GET("/users", func(c *gin.Context) {
        var users []User
        
        // Apply pagination to the query
        result := db.Scopes(pagination.Paginate(c)).Find(&users)
        if result.Error != nil {
            c.JSON(500, gin.H{"error": result.Error.Error()})
            return
        }
        
        // Get total count
        var total int64
        db.Model(&User{}).Count(&total)
        
        c.JSON(200, gin.H{
            "data": users,
            "page": c.Value(pagination.Page),
            "per_page": c.Value(pagination.Limit),
            "total": total,
        })
    })
    
    router.Run(":8080")
}
```

## Query Parameters
When making requests to your API endpoints, use these query parameters:
- `page`: The page number (default: 1)
- `per_page`: Number of items per page (default: 20, max: 100)

Example request:
```
GET /users?page=2&per_page=50
```

## Best Practices
1. Always apply the middleware at the router level or to specific route groups that need pagination
2. Handle the total count in your responses for proper client-side pagination
3. Consider implementing response metadata including:
   - Current page
   - Items per page
   - Total items
   - Total pages

## Notes
- The package is designed to work seamlessly with Gin and GORM
- All pagination parameters are validated and sanitized automatically
- The middleware is lightweight and adds minimal overhead
- The package uses context to pass pagination parameters between middleware and database utilities
