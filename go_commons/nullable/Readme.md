# Package nullable

## Overview
The nullable package provides a robust solution for handling nullable (optional) values in Go applications. Built on top of the `guregu/null` package, it offers type-safe wrappers and helper functions that prevent nil pointer dereferences and simplify working with optional data in your applications.

## Features
- Type-safe nullable types for common data types (String, Int, Float, Time, Bool)
- Automatic null handling based on zero values
- JSON marshaling/unmarshaling support
- Easy conversion between standard Go types and nullable types
- Prevention of nil pointer dereferences

## Available Types
The package provides wrapper functions for the following nullable types:
- `null.String`: For nullable string values
- `null.Int`: For nullable int64 values
- `null.Float`: For nullable float64 values
- `null.Time`: For nullable time.Time values
- `null.Bool`: For nullable boolean values

## Function Documentation

### NewNullableString
```go
func NewNullableString(s string) null.String
```
Creates a nullable string. Returns null if the string is empty, valid otherwise.

### NewNullableInt
```go
func NewNullableInt(i int64) null.Int
```
Creates a nullable int64. Returns null if the value is 0, valid otherwise.

### NewNullableFloat
```go
func NewNullableFloat(f float64) null.Float
```
Creates a nullable float64. Returns null if the value is 0.0, valid otherwise.

### NewNullableTime
```go
func NewNullableTime(t time.Time) null.Time
```
Creates a nullable time.Time. Returns null if the time is zero, valid otherwise.

### NewNullableBool
```go
func NewNullableBool(b bool) null.Bool
```
Creates a nullable boolean with the specified value.

## Usage Examples

### Working with Nullable Strings
```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/nullable"
)

func main() {
	// Creating valid nullable string
	validStr := nullable.NewNullableString("hello")
	fmt.Printf("Valid string: %v, IsValid: %v\n", validStr.String, validStr.Valid)

	// Creating null string
	nullStr := nullable.NewNullableString("")
	fmt.Printf("Null string: %v, IsValid: %v\n", nullStr.String, nullStr.Valid)
}
```

### Working with Nullable Numbers
```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/nullable"
)

func main() {
	// Integer example
	validInt := nullable.NewNullableInt(42)
	fmt.Printf("Valid int: %v, IsValid: %v\n", validInt.Int64, validInt.Valid)

	nullInt := nullable.NewNullableInt(0)
	fmt.Printf("Null int: %v, IsValid: %v\n", nullInt.Int64, nullInt.Valid)

	// Float example
	validFloat := nullable.NewNullableFloat(3.14)
	fmt.Printf("Valid float: %v, IsValid: %v\n", validFloat.Float64, validFloat.Valid)
}
```

### Working with Nullable Time
```go
package main

import (
	"fmt"
	"time"
	"github.com/omniful/go_commons/nullable"
)

func main() {
	// Valid time
	now := time.Now()
	validTime := nullable.NewNullableTime(now)
	fmt.Printf("Valid time: %v, IsValid: %v\n", validTime.Time, validTime.Valid)

	// Null time
	nullTime := nullable.NewNullableTime(time.Time{})
	fmt.Printf("Null time: %v, IsValid: %v\n", nullTime.Time, nullTime.Valid)
}
```

## Common Use Cases
1. Database interactions where fields might be NULL
2. JSON APIs with optional fields
3. Configuration settings with default values
4. Form handling with optional inputs

## Best Practices
1. Use nullable types for fields that can legitimately be NULL/undefined
2. Check the `Valid` field before accessing the value
3. Use the appropriate nullable type for each data type
4. Consider the zero-value behavior when creating nullable types

## Dependencies
- `gopkg.in/guregu/null.v4`: The underlying nullable types implementation

## Notes
- The package automatically handles NULL values based on Go's zero values
- All types support proper JSON marshaling/unmarshaling
- Thread-safe for concurrent operations
- Compatible with standard SQL database operations
