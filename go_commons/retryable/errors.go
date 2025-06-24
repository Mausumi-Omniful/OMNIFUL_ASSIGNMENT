package retryable

import "fmt"

// RetryExceededError represents an error when maximum retries have been exceeded
type RetryExceededError struct {
	// MaxRetries is the maximum number of retries that were attempted
	MaxRetries int
	// LastError is the last error that occurred before giving up
	LastError error
}

// Error implements the error interface
func (e *RetryExceededError) Error() string {
	return fmt.Sprintf("max retries (%d) exceeded: %v", e.MaxRetries, e.LastError)
}

// Unwrap returns the underlying error
func (e *RetryExceededError) Unwrap() error {
	return e.LastError
}

// IsRetryExceededError checks if the given error is a RetryExceededError
func IsRetryExceededError(err error) bool {
	_, ok := err.(*RetryExceededError)
	return ok
}
