package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
)

// LoggingMiddleware provides request/response logging using go_commons
func LoggingMiddleware() gin.HandlerFunc {
	return http.RequestLogMiddleware(http.LoggingMiddlewareOptions{
		Format:      "json",
		Level:       "info",
		LogRequest:  true,
		LogResponse: true,
		LogHeader:   false, // Set to false for security
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// The go_commons RequestLogMiddleware already handles request ID
		// This is a simple pass-through for now
		c.Next()
	}
} 