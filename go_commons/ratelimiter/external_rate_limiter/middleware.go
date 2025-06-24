package external_rate_limiter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	defaultLimit  = 10
	defaultExpiry = time.Minute
)

// ExternalRateLimitMiddleware is a gin middleware which limits requests per user using the provided rate limiter.
func ExternalRateLimitMiddleware(rateLimiter ExternalRateLimiter, options ...LimiterOption) gin.HandlerFunc {
	// Default options
	rlOptions := Options{
		Limit:  defaultLimit,
		Expiry: defaultExpiry,
	}

	// Apply custom options
	for _, opt := range options {
		opt(&rlOptions)
	}

	return func(c *gin.Context) {
		// Extract user identifier from request, such as IP address or user ID
		key := rlOptions.KeyExtractor(c)

		limit := rlOptions.Limit
		if rlOptions.LimitExtractor != nil {
			limit = rlOptions.LimitExtractor(c)
		}

		expiry := rlOptions.Expiry
		if rlOptions.ExpiryExtractor != nil {
			expiry = rlOptions.ExpiryExtractor(c)
		}

		// Perform rate limiting
		_, err := rateLimiter.IncrWithLimit(c, key, limit, expiry)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Request Quota exceeded ,Kindly slow down"})
			return
		}

		// Proceed to the next handler if the limit was not exceeded
		c.Next()
	}
}
