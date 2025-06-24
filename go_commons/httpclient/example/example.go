package example

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/httpclient"
	"github.com/omniful/go_commons/httpclient/auth"
	"github.com/omniful/go_commons/httpclient/request"
	"github.com/omniful/go_commons/httpclient/response"
	"github.com/omniful/go_commons/log"
	"net/url"
	"time"
)

// Global rate limiter to be used across requests
var rateLimiter = httpclient.NewRateLimiter(10, 1*time.Second)

type dummyAuthProvider struct{}

func (*dummyAuthProvider) Apply(_ *httpclient.Context, req request.Request) (request.Request, error) {
	// Do some process to apply custom auth
	req.GetHeaders().Add("Custom-Auth-Header", "random-value")
	return req, nil
}

func ExampleApi() {
	// Auth Mechanism
	a := auth.NewBasicAuth("foo", "bar")

	// Retry Strategy
	r := httpclient.NewLinearRetry(5*time.Second, 3)

	// Options
	opts := []httpclient.Option{
		httpclient.WithClientAuth(a),
		httpclient.WithRequestAuthProvider(&dummyAuthProvider{}),
		httpclient.WithRetry(r),
		httpclient.WithRateLimiter(rateLimiter),
		httpclient.WithLogConfig(httpclient.LogConfig{
			LogRequest:  true,
			LogResponse: true,
		}),

		// Callbacks
		httpclient.WithBeforeSendCallback(func(c *httpclient.Context, req request.Request) error {
			log.Info("Before Send Callback")
			return nil
		}),
		httpclient.WithBeforeAttemptCallback(func(c *httpclient.Context, req request.Request) error {
			log.Info("Before Attempt Callback")
			return nil
		}),
		httpclient.WithAfterAttemptCallback(func(c *httpclient.Context, req request.Request, resp response.Response) error {
			log.Info("After Attempt Callback")
			return nil
		}),
		httpclient.WithAfterSendCallback(func(c *httpclient.Context, req request.Request, resp response.Response) error {
			log.Info("After Send Callback")
			return nil
		}),
		httpclient.WithOnErrorCallback(func(c *httpclient.Context, req request.Request, resp response.Response, err error) {
			log.Info("On Error Callback")
		}),
	}

	// Api client
	c := httpclient.New("https://httpbin.org", opts...)

	// Prepare Request
	req, err := request.NewBuilder().
		SetUri("/{method}").
		SetHeaders(url.Values{
			"dummy-header": []string{"dummy-value"},
		}).
		SetPathParams(request.PathParams{
			"method": "post",
		}).
		SetBody(map[string]string{
			"foo": "bar",
		}).
		SetQueryParams(url.Values{
			"hello": []string{"world"},
		}).
		Build()
	if err != nil {
		panic(err)
	}

	// Send request
	resp, err := c.Post(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RAW RESPONSE BODY: %+v\n", resp)

	// Unmarshal Body
	var pb response.JsonBody
	err = resp.UnmarshalBody(&pb)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PARSED RESPONSE BODY: %+v\n", pb)
}
