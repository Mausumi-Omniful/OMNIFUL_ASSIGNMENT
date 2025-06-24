# Retryable Package

The **retryable** package provides a flexible and type-safe way to add retry logic to operations that may fail temporarily. It uses customizable options such as maximum retry attempts, exponential backoff with jitter & delay settings, and user-defined retry conditions. It supports both ordinary operations (returning only an error) and functions that return values.

The package is built upon three main components:

- **Backoff Calculation** (in `backoff.go`): Implements the exponential backoff algorithm with optional jitter.
- **Error Handling** (in `errors.go`): Defines error types such as `RetryExceededError` that are returned when all retry attempts fail.
- **Core Retry Logic** (in `retyable.go`): Provides the `Do` and `DoWithResult` functions along with a configurable `Config` type and several options to customize behavior.

## Features

- Generic error type support for type-safe retry conditions
- Configurable maximum retry attempts
- Configurable initial delay with exponential backoff
- Custom retry condition functions
- Retry attempt callbacks
- Context cancellation support
- Support for operations with and without return values

## Installation

```bash
go get github.com/omniful/go_commons
```

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/omniful/go_commons/retryable"
    "time"
)

// Simple retry with default configuration
err := retryable.Do(
    context.Background(),
    func() error {
        return someOperation()
    },
)

// Retry with custom configuration
err = retryable.Do(
    context.Background(),
    func() error {
        return someOperation()
    },
    retryable.WithMaxRetries[error](3),
    retryable.WithInitialDelay[error](100 * time.Millisecond),
)
```

### With Return Values

```go
// Retry operation that returns a value
result, err := retryable.DoWithResult(
    context.Background(),
    func() (string, error) {
        return someOperation()
    },
    retryable.WithMaxRetries[error](3),
)
```

### Custom Retry Conditions

```go
// Retry only on specific errors
err := retryable.Do(
    context.Background(),
    func() error {
        return someOperation()
    },
    retryable.WithIsRetryable(func(err error) bool {
        return errors.Is(err, io.ErrTemporary)
    }),
)
```

### With Retry Callbacks

```go
err := retryable.Do(
    context.Background(),
    func() error {
        return someOperation()
    },
    retryable.WithOnRetry(func(attempt int, err error) {
        log.Printf("Retry attempt %d after error: %v", attempt, err)
    }),
)
```

### Error Handling

```go
err := retryable.Do(
    context.Background(),
    someOperation,
)

if retryErr := new(retryable.RetryExceededError); errors.As(err, &retryErr) {
    fmt.Printf("Failed after %d retries. Last error: %v\n", 
        retryErr.MaxRetries, retryErr.LastError)
}

// Or use the helper function
if retryable.IsRetryExceededError(err) {
    // Handle retry exceeded case
}
```

## Configuration Options

### WithMaxRetries

Sets the maximum number of retry attempts.

```go
retryable.WithMaxRetries[error](3)
```

### WithInitialDelay

Sets the initial delay between retries. The delay increases exponentially with each retry.

```go
retryable.WithInitialDelay[error](100 * time.Millisecond)
```

### WithIsRetryable

Sets the function to determine if an error should be retried.

```go
retryable.WithIsRetryable(func(err error) bool {
    return errors.Is(err, io.ErrTemporary)
})
```

### WithOnRetry

Sets the function to be called before each retry attempt.

```go
retryable.WithOnRetry(func(attempt int, err error) {
    log.Printf("Retry attempt %d: %v", attempt, err)
})
```

## Error Types

### RetryExceededError

Returned when maximum retries have been exceeded.

```go
type RetryExceededError struct {
    MaxRetries int
    LastError  error
}
```

## Best Practices

1. **Context Usage**: Always pass a context with appropriate timeout to prevent infinite retry loops:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   
   err := retryable.Do(ctx, someOperation)
   ```

2. **Error Type Selection**: Choose appropriate error types for your use case:
   - Use standard `error` for simple cases
   - Use custom error types when you need specific retry conditions
   - Use error wrapping when you need to preserve error context

3. **Retry Conditions**: Carefully consider what errors should be retried:
   - Retry temporary failures (network issues, rate limits)
   - Don't retry permanent failures (validation errors, not found)
   - Consider using exponential backoff for rate limits

4. **Monitoring**: Use the OnRetry callback to monitor retry attempts:
   ```go
   retryable.WithOnRetry(func(attempt int, err error) {
       metrics.IncRetryCounter(attempt)
       log.Printf("Retry attempt %d: %v", attempt, err)
   })
   ```

## Advanced Usage

For advanced configuration such as custom error types, multiple error types, and fine-tuned retry conditions, please refer to the source code and additional examples.
