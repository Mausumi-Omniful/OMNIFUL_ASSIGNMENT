package response

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
)

type RedirectResponse struct {
	IsSuccess  bool        `json:"is_success"`
	StatusCode int         `json:"status_code"`
	Data       Redirect    `json:"data"`
	Meta       interface{} `json:"meta"`
}

type Redirect struct {
	Action   interface{} `json:"action"`
	Redirect interface{} `json:"redirect"`
	Payload  interface{} `json:"payload"`
}

func NewRedirectResponse(ctx *gin.Context, redirectData Redirect) {
	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: http.StatusMovedPermanently.Code(),
		Data:       redirectData,
	}

	ctx.AbortWithStatusJSON(http.StatusOK.Code(), res)
	return
}
