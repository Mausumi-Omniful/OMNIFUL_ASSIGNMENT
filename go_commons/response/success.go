package response

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
)

type SuccessResponse struct {
	IsSuccess  bool        `json:"is_success"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta"`
}

func NewSuccessResponse(ctx *gin.Context, data interface{}) {
	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: http.StatusOK.Code(),
		Data:       data,
	}

	ctx.AbortWithStatusJSON(http.StatusOK.Code(), res)
	return
}

func NewSuccessResponseWithMeta(ctx *gin.Context, data interface{}, meta interface{}) {
	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: http.StatusOK.Code(),
		Data:       data,
		Meta:       meta,
	}

	ctx.AbortWithStatusJSON(http.StatusOK.Code(), res)
	return
}

// NewSuccessResponseWithStatusCode creates a new success response with the specified HTTP status code,
// data, and metadata, and sends it as a JSON response using the provided Gin context.
//
// Parameters:
//   - ctx: The Gin context used to send the response.
//   - statusCode: The HTTP status code to be used for the response. If the provided status code
//     is not in the range of 200-299, it defaults to http.StatusOK.
//   - data: The payload data to include in the response.
//   - meta: Additional metadata to include in the response.
//
// Behavior:
//   - If the provided status code is not in the range of 200-299, it will be overridden with http.StatusOK.
//   - Constructs a SuccessResponse object with the provided data and metadata.
//   - Sends the response as JSON and aborts further request processing.
//
// Example:
//
//	// Example 1: Sending a success response with status code 201
//	data := map[string]string{"message": "Resource created successfully"}
//	meta := map[string]string{"request_id": "12345"}
//	NewSuccessResponseWithStatusCode(ctx, http.StatusCreated, data, meta)
//
//	// Example 2: Sending a success response with an invalid status code (defaults to 200)
//	data := map[string]string{"message": "Operation completed"}
//	NewSuccessResponseWithStatusCode(ctx, http.StatusBadRequest, data, nil)
func NewSuccessResponseWithStatusCode(ctx *gin.Context, statusCode http.StatusCode, data, meta any) {
	code := statusCode.Code()
	if code < 200 || code > 299 {
		code = http.StatusOK.Code()
	}

	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: code,
		Data:       data,
		Meta:       meta,
	}
	ctx.AbortWithStatusJSON(code, res)
}
