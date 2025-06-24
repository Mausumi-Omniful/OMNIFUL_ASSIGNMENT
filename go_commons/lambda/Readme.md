# Package lambda

## Overview
The lambda package provides a robust client for interacting with AWS Lambda functions in Go applications. It offers a clean, type-safe interface for executing Lambda functions, handling responses, and managing errors. The package is designed to simplify AWS Lambda integration while providing flexibility through configuration options.

## Features
- **AWS Lambda Client**: Easy-to-use client for invoking Lambda functions
- **Environment Prefixing**: Support for different environments (local, dev, staging, production)
- **Type-Safe Responses**: Structured response handling with automatic JSON deserialization
- **Error Handling**: Comprehensive error handling with detailed error responses
- **Mock Support**: Built-in mock client for testing
- **Configurable**: Flexible configuration through option pattern

## Installation
```go
go get github.com/omniful/go_commons
```

## Usage

### Creating a Lambda Client

1. Basic AWS Lambda Client:
```go
package main

import (
	"context"
	"github.com/omniful/go_commons/lambda"
)

func main() {
	ctx := context.Background()
	
	// Create a new AWS Lambda client with default configuration
	client, err := lambda.NewAWSLambdaClient(ctx)
	if err != nil {
		panic(err)
	}
}
```

2. Configured Lambda Client:
```go
// Create a client with custom environment prefix
client, err := lambda.NewAWSLambdaClient(ctx, 
	lambda.WithPrefix("staging"),
)
```

### Executing Lambda Functions

1. Basic Function Execution:
```go
type MyResponse struct {
	Message string `json:"message"`
}

func ExecuteLambda(ctx context.Context, client *lambda.Client) {
	var response MyResponse
	
	execReq := &request.ExecRequest{
		FunctionName: "my-function",
		Data: map[string]string{
			"key": "value",
		},
	}
	
	execRes, errRes := client.Execute(ctx, execReq, &response)
	if errRes != nil {
		// Handle error
		fmt.Printf("Error: %s - %s\n", errRes.ErrorCode, errRes.ErrorMessage)
		return
	}
	
	fmt.Printf("Success! Message: %s\n", response.Message)
}
```

2. Handling Lambda Responses:
```go
if execRes != nil {
	fmt.Printf("Status Code: %d\n", execRes.StatusCode)
	fmt.Printf("Executed Version: %s\n", execRes.ExecutedVersion)
	fmt.Printf("Response Data: %+v\n", execRes.Data)
}
```

## Response Format

### Success Response Structure
```json
{
	"status_code": 200,
	"data": {
		// Your response data here
	},
	"meta_data": {
		// Optional metadata
	}
}
```

### Error Response Structure
```json
{
	"status_code": 400,
	"data": {
		"error_code": "INVALID_PARAMS",
		"error_message": "Invalid parameters provided"
	},
	"meta_data": {
		// Optional error metadata
	}
}
```

## Testing
The package includes a mock client for testing purposes:

```go
func TestLambdaExecution(t *testing.T) {
	client, _ := lambda.NewMockClient("test")
	
	var response MyResponse
	execRes, errRes := client.Execute(context.TODO(),
		&request.ExecRequest{
			FunctionName: "test-function",
			Data: testData,
		},
		&response,
	)
	
	// Assert results
	assert.Nil(t, errRes)
	assert.NotNil(t, execRes)
}
```

## Best Practices
1. Always check for both `execRes` and `errRes` when executing Lambda functions
2. Use appropriate environment prefixes for different deployment stages
3. Implement proper error handling for both Lambda execution and response parsing
4. Use strongly typed structs for request and response data
5. Utilize the mock client for testing Lambda function interactions

## Notes
- The package is optimized for the stateless and ephemeral nature of AWS Lambda
- Local development can be configured using the `LOCAL_LAMBDA_ENDPOINT` environment variable
- Response handling automatically deserializes JSON responses into provided structs
- Error responses include detailed information about the failure including status code, error code, and message

## License
This package is part of the go_commons project. Please refer to the project's license file for terms of use.
