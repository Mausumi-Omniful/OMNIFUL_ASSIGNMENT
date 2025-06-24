# JSON Package

This package provides enhanced JSON marshaling and unmarshaling functionality with integrated New Relic error reporting. It wraps Go's standard `encoding/json` package and automatically reports any encoding/decoding errors to New Relic for better observability.

## Features

- Automatic error reporting to New Relic
- Context-aware error handling
- Request ID inclusion in error messages
- Drop-in replacement for standard json.Marshal/Unmarshal

## Installation

```go
import "github.com/omniful/go_commons/json"
```

## Usage

### Unmarshaling JSON

```go
package main

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/json"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	ctx := context.Background()
	
	// JSON data to unmarshal
	data := []byte(`{"name": "John Doe", "age": 30}`)
	
	var person Person
	err := json.Unmarshal(ctx, data, &person)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Person: %+v\n", person)
}
```

### Marshaling JSON

```go
package main

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/json"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	ctx := context.Background()
	
	person := Person{
		Name: "John Doe",
		Age:  30,
	}
	
	data, err := json.Marshal(ctx, person)
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}
	
	fmt.Printf("JSON: %s\n", string(data))
}
```

## Error Handling

Both `Marshal` and `Unmarshal` functions automatically report errors to New Relic with the following information:
- Request ID from context
- Priority level (P0)
- Detailed error message

## Dependencies

- `github.com/omniful/go_commons/env` - For request ID handling
- `github.com/omniful/go_commons/newrelic` - For error reporting

## Best Practices

1. Always pass a context with request ID when using these functions
2. Handle returned errors appropriately in your application
3. Monitor New Relic for any reported JSON encoding/decoding errors

## Note

This package is designed to be used in production environments where error monitoring is crucial. The automatic error reporting to New Relic helps in quickly identifying and debugging JSON-related issues.
