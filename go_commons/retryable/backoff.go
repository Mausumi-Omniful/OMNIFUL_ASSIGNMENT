package retryable

import (
	"math"
	"math/rand"
	"time"
)

// BackoffConfig defines the configuration for exponential backoff
type BackoffConfig struct {
	// InitialDelay is the initial delay between retries
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// Factor is the multiplier for exponential backoff
	Factor float64
	// Jitter enables randomization of the backoff duration
	Jitter bool
}

// DefaultBackoffConfig returns the default backoff configuration
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		InitialDelay: time.Second,
		MaxDelay:     30 * time.Second,
		Factor:       2.0,
		Jitter:       false,
	}
}

// calculateBackoff calculates the next backoff duration using exponential backoff algorithm
func calculateBackoff(attempt int, cfg BackoffConfig) time.Duration {
	// Calculate basic exponential backoff
	multiplier := math.Pow(cfg.Factor, float64(attempt))
	delay := float64(cfg.InitialDelay) * multiplier

	// Apply maximum delay if set
	if cfg.MaxDelay > 0 {
		delay = math.Min(delay, float64(cfg.MaxDelay))
	}

	// Add jitter if enabled
	if cfg.Jitter {
		// Add random jitter between 0% and 100% of the delay
		jitter := rand.Float64() // Random value between 0 and 1
		delay = delay * (1 + jitter)
	}

	return time.Duration(delay)
}
