# Package error

## Overview
The error package centralizes error handling, providing custom error types and helper functions to generate user-friendly error messages. This promotes consistency in error reporting across the codebase.

## Key Components
- Custom Errors: Structured errors that encapsulate error codes and messages.
- Error Messages: Functions to standardize error outputs.
- Error Types: Classified errors to differentiate between error conditions.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/error"
)

func main() {
	err := error.NewCustomError("example_error", "An example error occurred")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
~~~

## Notes
- Enhances debugging and user feedback through clear error definitions.
