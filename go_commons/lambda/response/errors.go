package response

import "github.com/omniful/go_commons/http"

// ErrorResponse is a error response returned by go commons lambda client
type ErrorResponse struct {
	StatusCode   http.StatusCode
	ErrorCode    string
	ErrorMessage string
	MetaData     map[string]interface{}
}

func NewErrorResponseBadRequest(errorCode string, errorMessage string) *ErrorResponse {
	return &ErrorResponse{
		StatusCode:   http.StatusBadRequest,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}
