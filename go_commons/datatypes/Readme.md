# Package datatypes

## Overview
The datatypes package defines common data structures and custom types used across the application, promoting consistency and reducing duplication.

## Key Components
- Custom Types: Structs and types representing key entities.
- Utility Functions: Operations to manipulate and validate these types.
- Centralization: Encourages reuse by centralizing type definitions.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/datatypes"
)

func main() {
	dt := datatypes.NewCustomType("example")
	fmt.Println("Data type instance:", dt)
}
~~~

## Notes
- Centralizes common type definitions to improve maintainability.
