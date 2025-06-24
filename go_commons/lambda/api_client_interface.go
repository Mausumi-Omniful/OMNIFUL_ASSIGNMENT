package lambda

import (
	"context"
	"github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/lambda/response"
)

type APIClientInterface interface {
	Invoke(context.Context, *request.InvokeRequest) (*response.InvokeResponse, error)
}
