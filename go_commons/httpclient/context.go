package httpclient

import (
	"context"
	"github.com/omniful/go_commons/httpclient/request"
	"time"
)

func NewContext(ctx context.Context, cfg Config, req request.Request) *Context {
	return &Context{ctx: ctx, cfg: cfg, req: req, attemptCount: 0}
}

type Context struct {
	ctx context.Context

	cfg          Config
	req          request.Request
	attemptCount int
}

/* Start context.Context interface implementation */
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

/* End context.Context interface implementation */

func (c *Context) Config() Config {
	return c.cfg
}

func (c *Context) Request() request.Request {
	return c.req
}

func (c *Context) AttemptCount() int {
	return c.attemptCount
}

func (c *Context) RecordAttempt() {
	c.attemptCount++
}
