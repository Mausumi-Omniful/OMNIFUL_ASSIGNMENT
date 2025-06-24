# Package pool

## Overview
The pool package provides mechanisms to manage a pool of reusable resources, such as connections or objects. It helps in reducing resource allocation overhead and improving performance.

## Key Components
- Resource Pool: Central data structure for resource management.
- Pooling Functions: Allocate, retrieve, and release resources.
- Configuration Options: Set pool sizes and timeouts.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/pool"
)

func main() {
	p := pool.New(10) // Pool with 10 resources
	resource := p.Get()
	// Use resource
	p.Put(resource)
	fmt.Println("Resource processed via pool.")
}
~~~

## Notes
- Proper configuration can greatly enhance application performance.
