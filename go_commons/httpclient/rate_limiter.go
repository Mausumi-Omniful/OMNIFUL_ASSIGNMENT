package httpclient

import (
	"go.uber.org/ratelimit"
	"time"
)

// TODO: This is not specific to httpclient. Can be extracted as generic rate limiter
type RateLimiter interface {
	WaitUntilReady()
}

func NewRateLimiter(noOfRequests int, per time.Duration) RateLimiter {
	l := ratelimit.New(noOfRequests, ratelimit.Per(per))
	return &rateLimiter{
		noOfRequests: noOfRequests,
		per:          per,
		limiter:      l,
	}
}

type rateLimiter struct {
	noOfRequests int
	per          time.Duration
	limiter      ratelimit.Limiter
}

func (r *rateLimiter) WaitUntilReady() {
	r.limiter.Take()
}
