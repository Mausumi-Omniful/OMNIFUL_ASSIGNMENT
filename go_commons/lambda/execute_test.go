package lambda

import (
	"context"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/lambda/lambda_api_client"
	"github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/lambda/response"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewMockClient(prefix string) (*Client, error) {
	client, err := NewLambdaClient(
		context.TODO(),
		WithPrefix(prefix),
		WithLambdaAPIClient(lambda_api_client.NewMockLambdaAPIClient()))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newMockReq(funcName string, expectedRes response.InvokeResponse) *request.ExecRequest {
	return &request.ExecRequest{
		FunctionName: funcName,
		Data:         expectedRes,
	}
}

func TestExecute(t *testing.T) {
	client, _ := NewMockClient("local")

	t.Run("LambdaSuccess2xx", func(t *testing.T) {
		t.Parallel()

		var res struct {
			TestKey string `json:"test_key"`
		}

		execRes, errRes := client.Execute(context.TODO(),
			newMockReq("mockSuccessFunc",
				response.InvokeResponse{
					ExecutedVersion: "testversion",
					StatusCode:      200,
					Payload:         []byte(`{"status_code":202,"data":{"test_key":"test_value"}}`),
				}), &res)

		assert.Nil(t, errRes, "error should be nil")

		assert.NotNil(t, execRes, "execRes shouldn't be nil")

		if execRes != nil {
			assert.Equal(t, http.StatusCode(202), execRes.StatusCode, "status code should be 202")
			assert.Equal(t, "testversion", execRes.ExecutedVersion, "executed version should be testversion")
			assert.Equal(t, "test_value", res.TestKey, "test_key value mismatch")
		}

	})

	t.Run("LambdaError4xx", func(t *testing.T) {
		t.Parallel()
		var res interface{}

		execRes, errRes := client.Execute(context.TODO(),
			newMockReq("mockErrorFunc",
				response.InvokeResponse{
					ExecutedVersion: "testversion",
					StatusCode:      200,
					Payload:         []byte(`{"status_code":400,"data":{"error_code":"INVALID_PARAMS","error_message":"Invalid Parameters"}}`),
				}), &res)

		assert.Nil(t, execRes, "execRes should be nil")

		assert.NotNil(t, errRes, "errorRes shouldn't be nil")
		if errRes != nil {
			assert.Equal(t, http.StatusCode(400), errRes.StatusCode, "status code should be 400")
			assert.Equal(t, "INVALID_PARAMS", errRes.ErrorCode, "error code not expected")
			assert.Equal(t, "Invalid Parameters", errRes.ErrorMessage, "error message not expected")
		}

	})

	t.Run("LambdaNonJsonFormatResponse", func(t *testing.T) {
		t.Parallel()
		var res interface{}

		execRes, errRes := client.Execute(context.TODO(),
			newMockReq("mockNonJsonResponseFunc",
				response.InvokeResponse{
					ExecutedVersion: "testversion",
					StatusCode:      200,
					Payload:         []byte(`<?xml version="1.0" encoding="UTF-8"?><errors><error><![CDATA[Login has already been taken]]></error></errors>`),
				}), &res)

		assert.Nil(t, execRes, "execRes should be nil")

		assert.NotNil(t, errRes, "errorRes shouldn't be nil")
		if errRes != nil {
			assert.Equal(t, http.StatusCode(400), errRes.StatusCode, "status code should be 400")
			assert.Equal(t, error2.JsonDeserializationError.ToString(), errRes.ErrorCode, "error code not expected")
		}
	})

}
