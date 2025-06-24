package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
)

func LoggingMiddleware() gin.HandlerFunc {
	return http.RequestLogMiddleware(http.LoggingMiddlewareOptions{
		Format:      "json",
		Level:       "info",
		LogRequest:  true,
		LogResponse: true,
		LogHeader:   false, 
	})
}

