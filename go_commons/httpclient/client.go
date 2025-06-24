package httpclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/httpclient/request"
	"github.com/omniful/go_commons/httpclient/response"
	"github.com/omniful/go_commons/log"
)

var ErrMaxRetriesExceeded = errors.New("max retries exceeded")

type Client interface {
	Get(context.Context, request.Request, ...Option) (response.Response, error)
	Head(context.Context, request.Request, ...Option) (response.Response, error)
	Options(context.Context, request.Request, ...Option) (response.Response, error)
	Post(context.Context, request.Request, ...Option) (response.Response, error)
	Put(context.Context, request.Request, ...Option) (response.Response, error)
	Patch(context.Context, request.Request, ...Option) (response.Response, error)
	Delete(context.Context, request.Request, ...Option) (response.Response, error)
	Send(context.Context, request.Request, ...Option) (response.Response, error)
}

func New(baseURl string, opts ...Option) Client {
	// Building config from options
	cfg := Options(opts).ToConfig()

	// Initialize resty client
	rc := resty.New().SetBaseURL(baseURl)

	// Check and configure client auth
	if a := cfg.clientAuth; a != nil {
		rc = rc.SetHeaders(ToAuthHeader(a))
	}

	// Check and set transport on resty client
	if t := cfg.transport; t != nil {
		rc = rc.SetTransport(newrelic.NewRoundTripper(t))
	} else {
		rc = rc.SetTransport(newrelic.NewRoundTripper(http.DefaultTransport))
	}

	// Check and set timeout
	if d := cfg.timeout; d > 0 {
		rc = rc.SetTimeout(d)
	}

	return &client{
		rc:  rc,
		cfg: cfg,
	}
}

type client struct {
	rc  *resty.Client
	cfg Config
}

func (c *client) Get(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodGet).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Head(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodHead).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Options(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodOptions).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Post(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodPost).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Put(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodPut).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Patch(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodPatch).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Delete(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	req, err := req.ToBuilder().SetMethod(http.MethodDelete).Build()
	if err != nil {
		return nil, err
	}
	return c.Send(ctx, req, opts...)
}

func (c *client) Send(ctx context.Context, req request.Request, opts ...Option) (response.Response, error) {
	// Initializing channels
	respch := make(chan response.Response)
	errch := NewSafeChannel[error]()

	defer close(respch)
	defer errch.Close()

	// Apply options on config
	cfg := c.cfg
	for _, opt := range opts {
		// apply config
		cfg = opt(cfg)
	}

	// Prepare and wrapping context in httpclient context
	hctx := c.prepareContext(ctx, req, cfg)

	// Apply timeout to the client if specified
	if cfg.timeout > 0 {
		c.rc = c.rc.SetTimeout(cfg.timeout)
	}

	go func() {
		// Recovery routine
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					// Wrap recover in error interface
					err = fmt.Errorf("%v", r)
				}
				if cfg.panicHandler != nil {
					cfg.panicHandler(err)
				}
				errch.Write(err) // propagate error to caller
			}
		}()

		resp, err := c.executeRequest(hctx, req, cfg)
		if err != nil {
			errch.Write(err)

			// Check and call on error Callback
			if cb := cfg.onError; cb != nil {
				executeOnErrorCallback(hctx, req, resp, err, cb)
			}

			return
		}
		respch <- resp
	}()

	select {
	case resp := <-respch:
		return resp, nil
	case err := <-errch.Read():
		return nil, err
	case <-hctx.Done():
		return nil, hctx.Err()
	}
}

func (c *client) executeRequest(ctx *Context, req request.Request, cfg Config) (resp response.Response, err error) {
	// Call before Send Callbacks
	if err := executeBeforeSendCallbacks(ctx, req, cfg.beforeSendCallbacks); err != nil {
		return nil, err
	}

	// Set request auth if request auth provider is set
	if ap := cfg.requestAuthProvider; ap != nil {
		req, err = ap.Apply(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	// Execute resty client
	resp, err = c.attemptRestyAndLog(ctx, req, cfg)
	if err != nil {
		return nil, err
	}

	retries := cfg.retryStrategies
	if len(retries) == 0 {
		// Retries are not configured
		return resp, err
	}

	// Retry is configured
	for resp.IsError() {
		//  Api failed. Identify retry strategy.
		retrych := make(chan Retry, 1)

		go func() {
			for _, retry := range retries {
				if retry.ShouldRetry(ctx, req, resp) {
					retrych <- retry
					break
				}
			}
			close(retrych)
		}()

		retry, ok := <-retrych
		if !ok {
			// No retry strategy identified. Breaking the loop
			break
		}

		// Retry strategy identified. Attempt to retry
		// Wait for next attempt time
		time.Sleep(retry.NextAttemptIn(ctx, req, resp))

		// Prepare next request for retry
		req, err = retry.PrepareRequest(ctx, req, resp)
		if err != nil {
			return nil, err
		}

		// Attempt retry
		resp, err = c.attemptRestyAndLog(ctx, req, cfg)
		if err != nil {
			return nil, err
		}
	}

	// Call after Send Callbacks
	if err := executeAfterSendCallbacks(ctx, req, resp, cfg.afterSendCallbacks); err != nil {
		return nil, err
	}

	return resp, err
}

func (c *client) attemptRestyAndLog(ctx *Context, req request.Request, cfg Config) (response.Response, error) {
	resp, err := c.attemptResty(ctx, req, cfg)
	if err != nil {
		c.logError(ctx, cfg, req, err)
		return nil, err
	}
	return resp, err
}

func (c *client) attemptResty(ctx *Context, req request.Request, cfg Config) (resp response.Response, err error) {
	if ctx.AttemptCount() >= cfg.maxRetries {
		return nil, ErrMaxRetriesExceeded
	}

	if rl := cfg.rateLimiter; rl != nil {
		// Rate limiter is configured. Wait until ready.
		rl.WaitUntilReady()
	}

	// Call before attempt Callbacks
	if err := executeBeforeAttemptCallbacks(ctx, req, cfg.beforeAttemptCallbacks); err != nil {
		return nil, err
	}

	// Record attempt in context
	ctx.RecordAttempt()

	// Execute resty request
	rs, err := c.prepareRestyRequest(ctx, req, cfg).Send()
	if err != nil {
		return nil, err
	}

	// Log successful attempt
	c.logRestyAttempt(ctx, cfg, rs)

	// Build response
	resp, err = response.NewResponse(c.rc, rs)
	if err != nil {
		return nil, err
	}

	// Call after attempt Callbacks
	if err := executeAfterAttemptCallbacks(ctx, req, resp, cfg.afterAttemptCallbacks); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) prepareContext(ctx context.Context, req request.Request, cfg Config) *Context {
	if cfg.forceRequestIDInHeaders {
		// Ensure request ID is present in headers
		if requestID := env.GetRequestID(ctx); len(requestID) == 0 {
			ctx = env.SetRequestID(ctx, env.NewRequestID())
		}
	}
	if cfg.deadline > 0 {
		ctx, _ = context.WithDeadline(ctx, time.Now().Add(cfg.deadline))
	}

	return NewContext(ctx, cfg, req)
}

func (c *client) prepareRestyRequest(ctx *Context, req request.Request, cfg Config) *resty.Request {
	r := c.rc.R().
		SetHeader(constants.HeaderUserAgent, cfg.userAgent).
		SetHeaderMultiValues(req.GetHeaders()).
		SetPathParams(req.GetPathParams()).
		SetQueryParamsFromValues(req.GetQueryParams()).
		SetBody(req.GetBody()).SetFormDataFromValues(req.GetFormData())
	r.URL = req.GetUri()
	r.Method = req.GetMethod()

	// Set request id in request headers
	if rID := env.GetRequestID(ctx); len(rID) > 0 {
		r = r.SetHeader(constants.HeaderXOmnifulRequestID, rID)
	}

	// Force content type if provided
	if len(cfg.contentType) > 0 {
		r = r.ForceContentType(cfg.contentType)
	}

	r.SetContext(ctx)

	return r
}

func (c *client) logRestyAttempt(ctx *Context, cfg Config, resp *resty.Response) {
	lc := cfg.logConfig
	if lc == nil {
		return
	}

	fields := map[string]any{
		"isSuccess":  resp.IsSuccess(),
		"httpStatus": resp.Status(),
	}
	if lc.LogRequest {
		req := resp.Request
		fields["request"] = map[string]any{
			"URL":         req.URL,
			"Method":      req.Method,
			"Header":      req.Header,
			"PathParams":  req.PathParams,
			"QueryParams": req.QueryParam,
			"FormData":    req.FormData,
			"Body":        req.Body,
		}
	}
	if lc.LogResponse {
		fields["response"] = map[string]any{
			"Body":   string(resp.Body()),
			"Header": resp.Header(),
		}
	}

	msg := "http call executed successfully"
	if !resp.IsSuccess() {
		msg = "http call failed"
	}

	c.log(ctx, cfg, fields, msg)
}

func (c *client) logError(ctx *Context, cfg Config, req request.Request, err error) {
	lc := cfg.logConfig
	if lc == nil {
		return
	}

	fields := map[string]any{
		"error": err,
	}
	if lc.LogRequest {
		fields["request"] = req
	}

	c.log(ctx, cfg, fields, fmt.Sprintf("http call failed with error: %v", err.Error()))
}

func (c *client) log(ctx *Context, cfg Config, fields map[string]any, msg string) {
	lc := cfg.logConfig
	if lc == nil {
		return
	}
	l := lc.Logger
	if l == nil {
		l = log.DefaultLogger()
	}
	l = l.WithFields(fields)

	lf := l.Info
	switch strings.ToLower(lc.LogLevel) {
	case "debug":
		lf = l.Debug
	case "info":
		lf = l.Info
	case "warn":
		lf = l.Warn
	case "error":
		lf = l.Error
	case "fatal":
		lf = l.Error
	case "panic":
		lf = l.Panic
	}

	logTag := fmt.Sprintf("[%s][httpclient] ", env.GetRequestID(ctx))
	lf(logTag + msg)
}
