package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
)

const (
	APIKeyHeader       = "X-API-Key"
	AuthorizationHeader = "Authorization"
)






// AuthMiddleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/test") || strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}

		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			authHeader := c.GetHeader(AuthorizationHeader)
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		expectedAPIKey := os.Getenv("OMS_API_KEY")
		if expectedAPIKey == "" {
			expectedAPIKey = "oms-dev-key-2025"
		}

		if apiKey == "" {
			fmt.Println("Missing API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.Translate(c.Request.Context(), "auth.missing_api_key"),
				"code":  "MISSING_API_KEY",
			})
			c.Abort()
			return
		}

		if apiKey != expectedAPIKey {
			fmt.Println("Invalid API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.Translate(c.Request.Context(), "auth.invalid_api_key"),
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey)
		fmt.Println("API key validated")
		c.Next()
	}
}




