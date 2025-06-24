package retryable

import (
	"context"
	"time"
)

// Config defines the retry configuration with generic error type
type Config[E error] struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// IsRetryable is a function that determines if an error should be retried
	IsRetryable func(E) bool
	// OnRetry is called before each retry attempt
	OnRetry func(attempt int, err E)
	// BackoffConfig defines the backoff strategy configuration
	BackoffConfig BackoffConfig
}

// DefaultConfig returns a default retry configuration
func DefaultConfig[E error]() Config[E] {
	return Config[E]{
		MaxRetries:    3,
		IsRetryable:   func(err E) bool { return any(err) != nil },
		OnRetry:       func(attempt int, err E) {},
		BackoffConfig: DefaultBackoffConfig(),
	}
}

// Option defines a function to modify Config
type Option[E error] func(*Config[E])

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries[E error](maxRetries int) Option[E] {
	return func(c *Config[E]) {
		c.MaxRetries = maxRetries
	}
}

// WithBackoffConfig sets the backoff configuration
func WithBackoffConfig[E error](cfg BackoffConfig) Option[E] {
	return func(c *Config[E]) {
		c.BackoffConfig = cfg
	}
}

// WithInitialDelay sets the initial delay between retries
func WithInitialDelay[E error](delay time.Duration) Option[E] {
	return func(c *Config[E]) {
		c.BackoffConfig.InitialDelay = delay
	}
}

// WithMaxDelay sets the maximum delay between retries
func WithMaxDelay[E error](delay time.Duration) Option[E] {
	return func(c *Config[E]) {
		c.BackoffConfig.MaxDelay = delay
	}
}

// WithJitter enables or disables jitter in the backoff calculation
func WithJitter[E error](enabled bool) Option[E] {
	return func(c *Config[E]) {
		c.BackoffConfig.Jitter = enabled
	}
}

// WithBackoffFactor sets the factor for exponential backoff
func WithBackoffFactor[E error](factor float64) Option[E] {
	return func(c *Config[E]) {
		c.BackoffConfig.Factor = factor
	}
}

// WithIsRetryable sets the function to determine if an error should be retried
func WithIsRetryable[E error](isRetryable func(E) bool) Option[E] {
	return func(c *Config[E]) {
		c.IsRetryable = isRetryable
	}
}

// WithOnRetry sets the function to be called before each retry attempt
func WithOnRetry[E error](onRetry func(attempt int, err E)) Option[E] {
	return func(c *Config[E]) {
		c.OnRetry = onRetry
	}
}

// Do executes the given function with retries according to the provided configuration
func Do[E error](ctx context.Context, fn func() E, opts ...Option[E]) error {
	// Apply default config
	cfg := DefaultConfig[E]()

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	var lastErr E
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute function
		err := fn()
		if any(err) == nil {
			return nil
		}

		lastErr = err

		// Check if we've reached max retries
		if attempt == cfg.MaxRetries {
			break
		}

		// Check if error is retryable
		if !cfg.IsRetryable(err) {
			return err
		}

		// Call OnRetry callback
		cfg.OnRetry(attempt+1, err)

		// Calculate and wait for backoff duration
		delay := calculateBackoff(attempt, cfg.BackoffConfig)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				// Drain the channel if timer already fired
				<-timer.C
			}
			return ctx.Err()
		case <-timer.C:
			// Timer has already fired and we've received from it
		}
	}

	return &RetryExceededError{
		MaxRetries: cfg.MaxRetries,
		LastError:  lastErr,
	}
}

// DoWithResult executes the given function with retries and returns its result
func DoWithResult[T any, E error](ctx context.Context, fn func() (T, E), opts ...Option[E]) (T, error) {
	var result T
	var lastResult T

	err := Do(ctx, func() E {
		var err E
		result, err = fn()
		if any(err) == nil {
			lastResult = result
		}
		return err
	}, opts...)

	return lastResult, err
}
