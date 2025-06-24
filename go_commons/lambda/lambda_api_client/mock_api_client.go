package lambda_api_client

import (
	"context"
	"github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/lambda/response"
)

// MockLambdaAPIClient is a mock implementation of the Lambda API Client for testing purposes.
type MockLambdaAPIClient struct{}

func (m *MockLambdaAPIClient) Invoke(ctx context.Context, req *request.InvokeRequest) (*response.InvokeResponse, error) {
	res := req.Data.(response.InvokeResponse)

	return &response.InvokeResponse{
		ExecutedVersion: res.ExecutedVersion,
		StatusCode:      res.StatusCode,
		Payload:         res.Payload,
	}, nil
}

func NewMockLambdaAPIClient() *MockLambdaAPIClient {
	return &MockLambdaAPIClient{}
}
