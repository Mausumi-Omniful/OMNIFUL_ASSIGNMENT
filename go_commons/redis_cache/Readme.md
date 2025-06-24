# Package redis_cache

## Overview
The redis_cache package offers caching functionalities built on Redis. It provides methods to set, get, and manage cache entries, thereby improving performance by reducing redundant computations.

## Key Components
- Cache Client: Manages connection to the Redis server.
- Caching Operations: Set, retrieve, and delete cache items.
- Configuration: Customizable options for cache behavior.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/redis_cache"
)

func main() {
	cache := redis_cache.NewClient(redis_cache.Config{
		// Configuration parameters
	})
	err := cache.Set("key", "value", 0)
	if err != nil {
		fmt.Println("Error setting cache:", err)
	} else {
		val, _ := cache.Get("key")
		fmt.Println("Cached value:", val)
	}
}
~~~

## Notes
- Integrates seamlessly with other Redis-based components.
