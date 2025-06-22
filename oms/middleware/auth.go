package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

const (
	// APIKeyHeader is the header name for API key authentication
	APIKeyHeader = "X-API-Key"
	// AuthorizationHeader is the standard authorization header
	AuthorizationHeader = "Authorization"
)

// AuthMiddleware provides API key authentication with i18n support
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for test endpoints and health check
		if strings.HasPrefix(c.Request.URL.Path, "/test") || strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}

		// Get API key from header
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			// Try Authorization header as fallback
			authHeader := c.GetHeader(AuthorizationHeader)
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Get expected API key from environment
		expectedAPIKey := os.Getenv("OMS_API_KEY")
		if expectedAPIKey == "" {
			// Use default for development
			expectedAPIKey = "oms-dev-key-2025"
		}

		// Validate API key
		if apiKey == "" {
			log.Warnf("❌ Missing API key for request: %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.Translate(c.Request.Context(), "auth.missing_api_key"),
				"code":  "MISSING_API_KEY",
			})
			c.Abort()
			return
		}

		if apiKey != expectedAPIKey {
			log.Warnf("❌ Invalid API key for request: %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": i18n.Translate(c.Request.Context(), "auth.invalid_api_key"),
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}

		// Add API key to context for potential use in handlers
		c.Set("api_key", apiKey)

		log.Infof("✅ API key validated for request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// OptionalAuthMiddleware provides optional API key authentication
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			// Try Authorization header as fallback
			authHeader := c.GetHeader(AuthorizationHeader)
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Get expected API key from environment
		expectedAPIKey := os.Getenv("OMS_API_KEY")
		if expectedAPIKey == "" {
			// Use default for development
			expectedAPIKey = "oms-dev-key-2025"
		}

		// If API key is provided, validate it
		if apiKey != "" {
			if apiKey != expectedAPIKey {
				log.Warnf("❌ Invalid API key for request: %s %s", c.Request.Method, c.Request.URL.Path)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": i18n.Translate(c.Request.Context(), "auth.invalid_api_key"),
					"code":  "INVALID_API_KEY",
				})
				c.Abort()
				return
			}
			c.Set("api_key", apiKey)
			c.Set("authenticated", true)
		} else {
			c.Set("authenticated", false)
		}

		c.Next()
	}
}
