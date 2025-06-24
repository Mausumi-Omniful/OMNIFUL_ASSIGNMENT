package response

import (
	"encoding/json"
	"fmt"
	errorCodes "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/http"
)

type InvokeResponse struct {
	ExecutedVersion string
	StatusCode      int
	Payload         []byte
}

func (invokeOutputRes *InvokeResponse) ToLambdaRes(resData interface{}) (*ExecResponse, *ErrorResponse) {
	lambdaStatusCode := http.StatusCode(invokeOutputRes.StatusCode)

	// Check if lambda invocation itself failed
	if !lambdaStatusCode.Is2xx() {
		return nil, &ErrorResponse{
			StatusCode:   lambdaStatusCode,
			ErrorCode:    fmt.Sprintf("%d", lambdaStatusCode),
			ErrorMessage: lambdaStatusCode.String(),
			MetaData: map[string]interface{}{
				"data": invokeOutputRes.Payload,
			},
		}
	}

	// Parse lambda response JSON
	lambdaRes := LambdaRes{}

	err := json.Unmarshal(invokeOutputRes.Payload, &lambdaRes)
	if err != nil {
		return nil, NewErrorResponseBadRequest(errorCodes.JsonDeserializationError.ToString(),
			err.Error())
	}

	statusCode := http.StatusCode(lambdaRes.StatusCode)

	// Handle success or error based on status code
	if statusCode.Is2xx() {
		if err = json.Unmarshal(lambdaRes.Data, resData); err != nil {
			return nil, NewErrorResponseBadRequest(errorCodes.JsonDeserializationError.ToString(),
				err.Error())
		}

		return &ExecResponse{
			ExecutedVersion: invokeOutputRes.ExecutedVersion,
			StatusCode:      statusCode,
			Data:            resData,
		}, nil
	}

	// Handle error responses
	errorData := ErrorLambdaData{}
	if err = json.Unmarshal(lambdaRes.Data, &errorData); err != nil {
		return nil, NewErrorResponseBadRequest(errorCodes.JsonDeserializationError.ToString(),
			err.Error())
	}

	return nil, &ErrorResponse{
		StatusCode:   statusCode,
		ErrorCode:    errorData.ErrorCode,
		ErrorMessage: errorData.ErrorMessage,
		MetaData:     lambdaRes.MetaData,
	}
}
