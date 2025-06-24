package interservice_client

import (
	"context"
	"encoding/json"
	"errors"
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"github.com/omniful/go_commons/response"
	"net"
)

type Client struct {
	httpClient *http.Client
}

type InterSvcResponse struct {
	IsSuccess  bool                   `json:"is_success"`
	StatusCode int                    `json:"status_code"`
	Data       json.RawMessage        `json:"data"`
	Meta       map[string]interface{} `json:"meta"`
	Error      json.RawMessage        `json:"error"`
}

type Error struct {
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors"`
	StatusCode http.StatusCode   `json:"status_code"`
	Data       interface{}       `json:"data"`
}

type ErrorWithData struct {
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors"`
	Data       interface{}       `json:"data"`
	StatusCode http.StatusCode   `json:"status_code"`
}

func NewClient(client *http.Client) *Client {
	return &Client{
		httpClient: client,
	}
}

func (c *Client) Get(request *http.Request, data interface{}) (*resty.Response, *Error) {
	return c.execute(context.Background(), http.APIGet, request, data)
}

func (c *Client) Post(request *http.Request, data interface{}) (*resty.Response, *Error) {
	return c.execute(context.Background(), http.APIPost, request, data)
}

func (c *Client) Put(request *http.Request, data interface{}) (*resty.Response, *Error) {
	return c.execute(context.Background(), http.APIPut, request, data)
}

func (c *Client) Patch(request *http.Request, data interface{}) (*resty.Response, *Error) {
	return c.execute(context.Background(), http.APIPatch, request, data)
}

func (c *Client) Delete(request *http.Request, data interface{}) (*resty.Response, *Error) {
	return c.execute(context.Background(), http.APIDelete, request, data)
}

func (c *Client) execute(
	ctx context.Context,
	method http.APIMethod,
	request *http.Request,
	data interface{},
) (*resty.Response, *Error) {
	result := InterSvcResponse{}
	res, err := c.httpClient.Execute(ctx, method, request, &result)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"tag":          "interservice response",
			"status_code":  res.StatusCode(),
			"httpResponse": string(res.Body()),
		}).Error("Error executing request")
		return nil, &Error{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}

	return parseResponse(ctx, &result, data, res)
}

func (c *Client) Execute(
	ctx context.Context,
	method http.APIMethod,
	request *http.Request,
	data interface{},
) (*InterSvcResponse, *Error) {
	result := InterSvcResponse{}
	res, err := c.httpClient.Execute(ctx, method, request, &result)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"tag":                             "Interservice-response",
			"error_message":                   err.Error(),
			"status_code":                     res.StatusCode(),
			"httpResponse":                    res.String(),
			constants.HeaderXOmnifulRequestID: env.GetRequestID(ctx),
		}).Errorf("Error executing request")

		// Handle Context Cancelled
		if errors.Is(err, context.Canceled) {
			return nil, &Error{
				Message:    http.StatusRequestTimeout.String(),
				StatusCode: http.StatusRequestTimeout}
		}

		// Handle Deadline Exceeded
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &Error{
				Message:    http.StatusRequestTimeout.String(),
				StatusCode: http.StatusRequestTimeout}
		}

		// Handle Timeout Errors
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, &Error{
				Message:    http.StatusRequestTimeout.String(),
				StatusCode: http.StatusRequestTimeout}
		}

		// For Others, Internal Server Errors
		return nil, &Error{Message: http.StatusInternalServerError.String(), StatusCode: http.StatusInternalServerError}
	}

	_, parseErr := parseResponse(ctx, &result, data, res)
	if parseErr != nil {
		return nil, parseErr
	}

	return &result, nil
}

func (c *Client) ExecuteWithValidation(
	ctx context.Context,
	method http.APIMethod,
	request *http.Request,
	data interface{},
	validator *validatorPkg.Validate,
) (*InterSvcResponse, *Error) {
	result, interSvcErr := c.Execute(ctx, method, request, data)
	if interSvcErr != nil {
		return result, interSvcErr
	}

	//Validator segment
	segment := newrelic.StartSegmentWithContext(ctx, "Validator")
	defer segment.End()

	err := validator.VarCtx(ctx, data, "dive")
	if err != nil {
		return nil, &Error{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}

	return result, nil
}

func parseResponse(ctx context.Context,
	result *InterSvcResponse, data interface{}, httpResponse *resty.Response) (*resty.Response, *Error) {
	// For success 2xx Codes: unmarshall response body to predefined success response
	if httpResponse.IsSuccess() {
		err := json.Unmarshal(result.Data, data)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"tag":                             "Interservice-response",
				"error_message":                   err.Error(),
				"status_code":                     httpResponse.StatusCode(),
				"httpResponse":                    httpResponse.String(),
				constants.HeaderXOmnifulRequestID: env.GetRequestID(ctx),
			}).Errorf("Unable to unmarshal success response : %+v", result.Data)

			return httpResponse, &Error{Message: http.StatusInternalServerError.String(), StatusCode: http.StatusInternalServerError}
		}

		return httpResponse, nil
	}

	// For Error, 4XX Codes: unmarshall error response to predefined error response
	if httpResponse.StatusCode() > 399 && httpResponse.StatusCode() < 500 {
		var resErr response.Error

		err := json.Unmarshal(result.Error, &resErr)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"tag":                             "Interservice-response",
				"error_message":                   err.Error(),
				"status_code":                     httpResponse.StatusCode(),
				"httpResponse":                    httpResponse.String(),
				constants.HeaderXOmnifulRequestID: env.GetRequestID(ctx),
			}).Errorf("Unable to unmarshal error response : %+v", result.Error)

			return nil, &Error{Message: http.StatusInternalServerError.String(), StatusCode: http.StatusInternalServerError}
		}

		return nil, &Error{
			Message:    resErr.Message,
			Errors:     resErr.Errors,
			StatusCode: http.StatusCode(httpResponse.StatusCode()),
			Data:       resErr.Data,
		}
	}

	// Rest Status Codes (1xx,3xx,5xx)
	statusCode := http.StatusCode(httpResponse.StatusCode())

	return nil, &Error{
		Message:    statusCode.String(),
		StatusCode: statusCode,
	}
}
