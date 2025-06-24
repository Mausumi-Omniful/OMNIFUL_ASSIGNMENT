package interservice_client

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/response"
)

func NewErrorResponseByInterServiceError(ctx *gin.Context, error *Error) {
	res := &response.ErrorResponse{
		IsSuccess:  false,
		StatusCode: error.StatusCode.Code(),
		Error: response.Error{
			Message: error.Message,
			Errors:  error.Errors,
		},
	}

	ctx.AbortWithStatusJSON(error.StatusCode.Code(), res)
	return
}

func NewErrorWithDataResponseByInterServiceError(ctx *gin.Context, error *Error) {
	res := &response.ErrorResponse{
		IsSuccess:  false,
		StatusCode: error.StatusCode.Code(),
		Error: response.Error{
			Message: error.Message,
			Errors:  error.Errors,
			Data:    error.Data,
		},
	}

	ctx.AbortWithStatusJSON(error.StatusCode.Code(), res)
	return
}
