package file_utilities

import (
	"context"
	commonErr "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/file_utilities/requests"
	"github.com/omniful/go_commons/file_utilities/responses"
	"github.com/omniful/go_commons/lambda"
	lambdaReq "github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/log"
)

// GenerateAWBLabel:
// Triggers an AWS Lambda function to generate an Air Waybill (AWB) label PDF based on the provided shipment data and configuration.
//
// Note: Maximum of 100 data entries supported in data array.
func GenerateAWBLabel(
	ctx context.Context,
	lambdaClient *lambda.Client,
	printAWBLabelRequest requests.PrintAWBLabelRequest,
) (responses.PrintAWBLabelResponse, commonErr.CustomError) {

	// validate the request
	cusErr := validatePrintAWBLabelRequest(ctx, lambdaClient, printAWBLabelRequest)
	if cusErr.Exists() {
		return responses.PrintAWBLabelResponse{}, cusErr
	}

	res := responses.PrintAWBLabelResponse{}

	_, errorRes := lambdaClient.Execute(ctx, &lambdaReq.ExecRequest{
		FunctionName: GenerateAWBLabelLambda,
		Data:         printAWBLabelRequest,
	}, &res)
	if errorRes != nil {
		return responses.PrintAWBLabelResponse{}, commonErr.NewCustomError(commonErr.BadRequestError, errorRes.ErrorMessage)
	}

	return res, commonErr.CustomError{}
}

// validatePrintAWBLabelRequest:
// Validates the input request for AWB label generation to ensure all required fields are present and correctly formatted before invoking the Lambda function.
// Prevents unnecessary Lambda invocations due to malformed input.
func validatePrintAWBLabelRequest(
	ctx context.Context,
	lambdaClient *lambda.Client,
	req requests.PrintAWBLabelRequest,
) commonErr.CustomError {
	if lambdaClient == nil {
		log.ErrorfWithContext(ctx, "lambdaClient is nil")
		return commonErr.NewCustomError(commonErr.RequestNotValid, "lambda client is nil")
	}

	if len(req.InputTemplateFileName) == 0 {
		log.ErrorfWithContext(ctx, "template file name is empty")
		return commonErr.NewCustomError(commonErr.RequestNotValid, "template file name is empty")
	}

	if len(req.Service) == 0 {
		log.ErrorfWithContext(ctx, "service name is empty")
		return commonErr.NewCustomError(commonErr.RequestNotValid, "service name is empty")
	}

	// Check if req.Data is empty (covers nil as well)
	if len(req.Data) == 0 {
		log.ErrorfWithContext(ctx, "request data is empty")
		return commonErr.NewCustomError(commonErr.RequestNotValid, "request data cannot be empty")
	}

	return commonErr.CustomError{}
}
