package file_utilities

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/env"
	commonErr "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/file_utilities/requests"
	"github.com/omniful/go_commons/file_utilities/responses"
	"github.com/omniful/go_commons/lambda"
	lambdaReq "github.com/omniful/go_commons/lambda/request"
	"github.com/omniful/go_commons/log"
)

func PdfMerger(
	ctx context.Context,
	lambdaClient *lambda.Client,
	pdfMergeRequest requests.PdfMergeRequest,
) (responses.PdfMergeResponse, commonErr.CustomError) {
	logTag := fmt.Sprintf("RequestID: %s Function: PdfMerger ", env.GetRequestID(ctx))

	if lambdaClient == nil {
		log.Errorf(logTag + "lambdaClient is nil")

		return responses.PdfMergeResponse{}, commonErr.NewCustomError(commonErr.RequestNotValid, "lambda client is nil")
	}

	if len(pdfMergeRequest.PdfURLs) == 0 {
		log.Errorf(logTag + "empty pdf urls provided")

		return responses.PdfMergeResponse{}, commonErr.NewCustomError(commonErr.RequestNotValid, "pdf urls not provided")
	}

	res := responses.PdfMergeResponse{}

	_, errorRes := lambdaClient.Execute(ctx, &lambdaReq.ExecRequest{
		FunctionName: MergePDFLambda,
		Data:         pdfMergeRequest,
	}, &res)
	if errorRes != nil {
		return responses.PdfMergeResponse{}, commonErr.NewCustomError(commonErr.BadRequestError, errorRes.ErrorMessage)
	}

	return res, commonErr.CustomError{}
}
