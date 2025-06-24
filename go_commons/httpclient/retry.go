package httpclient

import (
	"github.com/omniful/go_commons/httpclient/request"
	"github.com/omniful/go_commons/httpclient/response"
	"math"
	"time"
)

type Retry interface {
	ShouldRetry(*Context, request.Request, response.Response) bool
	NextAttemptIn(*Context, request.Request, response.Response) time.Duration
	PrepareRequest(*Context, request.Request, response.Response) (request.Request, error)
}

/* START FIXED RETRY */

func NewFixedRetry(backoff time.Duration, maxAttempts int) Retry {
	return &fixedRetry{backoff: backoff, maxAttempts: maxAttempts}
}

type fixedRetry struct {
	backoff     time.Duration
	maxAttempts int
}

func (r fixedRetry) ShouldRetry(ctx *Context, _ request.Request, _ response.Response) bool {
	return ctx.AttemptCount() < r.maxAttempts
}

func (r fixedRetry) NextAttemptIn(_ *Context, _ request.Request, _ response.Response) time.Duration {
	return r.backoff
}

func (r fixedRetry) PrepareRequest(_ *Context, req request.Request, _ response.Response) (request.Request, error) {
	return req, nil
}

/* END FIXED RETRY */

/* START LINEAR RETRY */

func NewLinearRetry(backoff time.Duration, maxAttempts int) Retry {
	return &linearRetry{backoff: backoff, maxAttempts: maxAttempts}
}

type linearRetry struct {
	backoff     time.Duration
	maxAttempts int
}

func (r linearRetry) ShouldRetry(ctx *Context, _ request.Request, _ response.Response) bool {
	return ctx.AttemptCount() < r.maxAttempts
}

func (r linearRetry) NextAttemptIn(ctx *Context, _ request.Request, _ response.Response) time.Duration {
	return time.Duration(int64(ctx.AttemptCount()) * r.backoff.Nanoseconds())
}

func (r linearRetry) PrepareRequest(_ *Context, req request.Request, _ response.Response) (request.Request, error) {
	return req, nil
}

/* END LINEAR RETRY */

/* START EXPONENTIAL RETRY */

func NewExponentialRetry(delay time.Duration, maxAttempts int) Retry {
	return &exponentialRetry{delay: delay, maxAttempts: maxAttempts}
}

type exponentialRetry struct {
	delay       time.Duration
	maxAttempts int
}

func (r exponentialRetry) ShouldRetry(ctx *Context, _ request.Request, _ response.Response) bool {
	return ctx.AttemptCount() < r.maxAttempts
}

func (r exponentialRetry) NextAttemptIn(ctx *Context, _ request.Request, _ response.Response) time.Duration {
	return time.Duration(int64(math.Pow(2, float64(ctx.AttemptCount()))) * r.delay.Nanoseconds())
}

func (r exponentialRetry) PrepareRequest(_ *Context, req request.Request, _ response.Response) (request.Request, error) {
	return req, nil
}

/* END EXPONENTIAL RETRY */
