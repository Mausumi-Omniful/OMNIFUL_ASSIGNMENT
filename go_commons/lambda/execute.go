package lambda

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/env"
	errorCodes "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/lambda/response"
	"github.com/omniful/go_commons/log"
)

// Execute calls the lambda function and returns a successful Response if the lambda completes with a 2xx status code.
// It also unmarshal the response data into the provided response struct.
//
// If an error occurs, Execute returns nil for Response, skips unmarshalling, and instead provides an ErrorResponse.
//
// Lambda Response Conventions:
// - Responses should always be in JSON format.
// - Successful response format: {"status_code":"200","data":{}}
// - Error response format: {"status_code":"400","data":{"error_code":"INVALID_PARAMS","error_message":"Invalid parameters"}}
// - For any unexpected response structure, Execute returns an ErrorResponse with Status 400 and error code JSON_DESERIALIZATION_ERROR.
func (c *Client) Execute(ctx context.Context, req *request.ExecRequest, res interface{}) (*response.ExecResponse, *response.ErrorResponse) {
	logTag := fmt.Sprintf("RequestID: %s Function: LambdaExecute ", env.GetRequestID(ctx))

	if c == nil {
		log.Errorf(logTag + "lambda client is nil")

		return nil,
			response.NewErrorResponseBadRequest(errorCodes.BadRequestError.ToString(), "lambda client is nil")
	}

	if len(req.FunctionName) == 0 {
		log.Errorf(logTag + "no function name given")

		return nil,
			response.NewErrorResponseBadRequest(errorCodes.BadRequestError.ToString(), "function name is empty")
	}

	invokeRes, err := c.lambdaAPIClient.Invoke(ctx,
		&request.InvokeRequest{
			ExecRequest: req,
			Namespace:   c.prefix,
		})
	if err != nil {
		log.Errorf(logTag+"lambda invocation failed,err: %+v", err.Error())

		return nil, response.NewErrorResponseBadRequest(errorCodes.BadRequestError.ToString(), err.Error())
	}

	lambdaRes, errRes := invokeRes.ToLambdaRes(res)
	if errRes != nil {
		log.Errorf(logTag+"failed to convert lambda res,errRes: %+v", errRes)

		return nil, errRes
	}

	return lambdaRes, nil
}
