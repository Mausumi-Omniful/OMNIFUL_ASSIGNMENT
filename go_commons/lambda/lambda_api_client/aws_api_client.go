package lambda_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsLambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	awsLambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/newrelic/go-agent/v3/integrations/nrawssdk-v2"
	"github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/lambda/response"
)

type AWSLambdaAPIClient struct {
	*awsLambda.Client
}

func NewAWSLambdaAPIClient(ctx context.Context) (*AWSLambdaAPIClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, func(awsConfig *config.LoadOptions) error {
		// Instrument all new AWS clients with New Relic
		nrawssdk.AppendMiddlewares(&awsConfig.APIOptions, nil)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	opts := make([]func(options *awsLambda.Options), 0)

	if endpoint, ok := os.LookupEnv("LOCAL_LAMBDA_ENDPOINT"); ok {
		opts = append(opts, func(o *awsLambda.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	awsClient := awsLambda.NewFromConfig(cfg, opts...)

	return &AWSLambdaAPIClient{awsClient}, nil
}

func (c *AWSLambdaAPIClient) Invoke(ctx context.Context, request *request.InvokeRequest) (*response.InvokeResponse, error) {
	payload, err := json.Marshal(request.Data)
	if err != nil {
		return nil,
			fmt.Errorf("request data not json serializable")
	}

	lambdaInvokeRes, err := c.Client.Invoke(ctx, &awsLambda.InvokeInput{
		FunctionName:   aws.String(request.Namespace + "_" + request.FunctionName),
		InvocationType: awsLambdaTypes.InvocationTypeRequestResponse,
		LogType:        awsLambdaTypes.LogTypeNone,
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	return &response.InvokeResponse{
		ExecutedVersion: *lambdaInvokeRes.ExecutedVersion,
		StatusCode:      int(lambdaInvokeRes.StatusCode),
		Payload:         lambdaInvokeRes.Payload,
	}, nil
}
