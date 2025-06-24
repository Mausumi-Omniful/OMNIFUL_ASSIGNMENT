package external_rate_limiter

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// ExternalRateLimiter is an interface for rate limiting mechanisms applied to external requests.
type ExternalRateLimiter interface {
	IncrWithLimit(ctx context.Context, key string, limit int, expiry time.Duration) (int64, error)
}

// Options represents options for ExternalRateLimitMiddleware.
type Options struct {
	Limit           int
	Expiry          time.Duration
	KeyExtractor    func(*gin.Context) string        // Function to extract the key from the request
	LimitExtractor  func(*gin.Context) int           // Function to extract the limit from the request
	ExpiryExtractor func(*gin.Context) time.Duration // Function to extract the expiry from the request
}

// LimiterOption is a functional option for ExternalRateLimitMiddleware.
type LimiterOption func(*Options)

// WithLimit sets the rate limit.
func WithLimit(limit int) LimiterOption {
	return func(o *Options) {
		o.Limit = limit
	}
}

// WithExpiry sets the expiry duration.
func WithExpiry(expiry time.Duration) LimiterOption {
	return func(o *Options) {
		o.Expiry = expiry
	}
}

// WithKeyExtractor sets the key extractor function.
func WithKeyExtractor(extractor func(*gin.Context) string) LimiterOption {
	return func(o *Options) {
		o.KeyExtractor = extractor
	}
}

func WithLimitExtractor(extractor func(*gin.Context) int) LimiterOption {
	return func(o *Options) {
		o.LimitExtractor = extractor
	}
}

func WithExpiryExtractor(extractor func(*gin.Context) time.Duration) LimiterOption {
	return func(o *Options) {
		o.ExpiryExtractor = extractor
	}
}
