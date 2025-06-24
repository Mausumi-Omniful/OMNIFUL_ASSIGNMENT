package lambda

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/lambda/lambda_api_client"
)

type Client struct {
	prefix          string // Used to add prefix to differentiate between running environments
	lambdaAPIClient APIClientInterface
}

type Option func(client *Client) error

// NewLambdaClient initializes a Lambda client with optional configurations.
// Allows setting environment and additional options via Option functions.
// Future Scopes: Add new relic instrumentation to monitor lambda performance
func NewLambdaClient(ctx context.Context, opts ...Option) (*Client, error) {
	client := &Client{
		prefix:          "local", // Default
		lambdaAPIClient: nil,
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	if client.lambdaAPIClient == nil {
		return nil, fmt.Errorf("lambda api client is nil")
	}

	return client, nil
}

// NewAWSLambdaClient initializes an AWS Lambda client with optional configurations.
// Allows setting environment and additional options via Option functions.
func NewAWSLambdaClient(ctx context.Context, opts ...Option) (*Client, error) {
	awsLambdaAPIClient, err := lambda_api_client.NewAWSLambdaAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	opts = append(opts, WithLambdaAPIClient(awsLambdaAPIClient))

	client, err := NewLambdaClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// WithPrefix configures the Lambda client for a specified environment (local, dev, staging, production ...).
func WithPrefix(prefix string) Option {
	return func(client *Client) error {
		client.prefix = prefix

		return nil
	}
}

// WithLambdaAPIClient configures the Lambda client to specify api client for lambda invocation.
func WithLambdaAPIClient(apiClientInterface APIClientInterface) Option {
	return func(client *Client) error {
		client.lambdaAPIClient = apiClientInterface

		return nil
	}
}

func (c *Client) GetPrefix() string {
	return c.prefix
}
