package response

import (
	"sync"

	"github.com/gin-gonic/gin"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

type ErrorResponse struct {
	IsSuccess  bool  `json:"is_success"`
	StatusCode int   `json:"status_code"`
	Error      Error `json:"error"`
}

type Error struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

var customCodeToHttpCodeMapping = map[error2.Code]http.StatusCode{}

var once sync.Once

// Set custom error mapping
func SetCustomErrorMapping(mapping map[error2.Code]http.StatusCode) {
	once.Do(func() {
		customCodeToHttpCodeMapping = mapping
	})
}

type Option func(*ErrorResponse)

func NewErrorResponse(
	ctx *gin.Context,
	customError error2.CustomError,
	customErrorCodeToErrorRespMapping map[error2.Code]http.StatusCode,
	options ...Option,
) {
	statusCode, ok := customErrorCodeToErrorRespMapping[customError.ErrorCode()]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	generateErrorResponse(ctx, customError, statusCode, options...)
}

func NewErrorResponseByStatusCode(ctx *gin.Context, statusCode http.StatusCode) {
	res := &ErrorResponse{
		IsSuccess:  false,
		StatusCode: statusCode.Code(),
		Error: Error{
			Message: statusCode.String(),
		},
	}

	ctx.AbortWithStatusJSON(statusCode.Code(), res)
}

func NewErrorResponseV2(
	ctx *gin.Context,
	customError error2.CustomError,
	options ...Option,
) {
	statusCode, ok := customCodeToHttpCodeMapping[customError.ErrorCode()]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	generateErrorResponse(ctx, customError, statusCode, options...)
}

func generateErrorResponse(
	ctx *gin.Context,
	customError error2.CustomError,
	statusCode http.StatusCode,
	options ...Option,
) {

	if customError.Exists() && statusCode != http.StatusOK {
		log.WithFields(map[string]interface{}{
			"host":     ctx.Request.Host,
			"path":     ctx.Request.URL.Path,
			"query":    ctx.Request.URL.RawQuery,
			"method":   ctx.Request.Method,
			"is_error": true,
			"error":    customError,
		}).Errorf("APIFailed with err: %s", customError.Error())
	}

	var message string
	if customError.ErrorCode() == error2.RequestInvalid {
		message = customError.UserMessage()
	} else {
		message = statusCode.String()
	}

	res := &ErrorResponse{
		IsSuccess:  false,
		StatusCode: statusCode.Code(),
		Error: Error{
			Message: message,
			Data:    customError.ErrorData(),
			Errors:  customError.ErrorMap(),
		},
	}

	for _, option := range options {
		option(res)
	}

	ctx.AbortWithStatusJSON(statusCode.Code(), res)
}
