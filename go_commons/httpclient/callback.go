package httpclient

import (
	"github.com/omniful/go_commons/httpclient/request"
	"github.com/omniful/go_commons/httpclient/response"
)

type BeforeCallback func(*Context, request.Request) error
type AfterCallback func(*Context, request.Request, response.Response) error
type OnErrorCallback func(*Context, request.Request, response.Response, error)

type BeforeSendCallback BeforeCallback
type BeforeAttemptCallback BeforeCallback
type AfterAttemptCallback AfterCallback
type AfterSendCallback AfterCallback

func executeBeforeSendCallbacks(ctx *Context, req request.Request, cbs []BeforeSendCallback) error {
	for _, cb := range cbs {
		if err := cb(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func executeBeforeAttemptCallbacks(ctx *Context, req request.Request, cbs []BeforeAttemptCallback) error {
	for _, cb := range cbs {
		if err := cb(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func executeAfterAttemptCallbacks(ctx *Context, req request.Request, resp response.Response, cbs []AfterAttemptCallback) error {
	for _, cb := range cbs {
		if err := cb(ctx, req, resp); err != nil {
			return err
		}
	}
	return nil
}

func executeAfterSendCallbacks(ctx *Context, req request.Request, resp response.Response, cbs []AfterSendCallback) error {
	for _, cb := range cbs {
		if err := cb(ctx, req, resp); err != nil {
			return err
		}
	}
	return nil
}

func executeOnErrorCallback(ctx *Context, req request.Request, resp response.Response, err error, cb OnErrorCallback) {
	// ensure resp is not nil
	if resp == nil {
		resp = response.EmptyResponse()
	}

	cb(ctx, req, resp, err)
}
