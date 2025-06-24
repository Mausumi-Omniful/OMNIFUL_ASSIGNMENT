package env

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/constants"
)

// RequestID checks the X-Request-ID header and generate new if not found
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if correlationID := c.GetHeader(constants.HeaderXOmnifulCorrelationID); len(correlationID) > 0 {
			c.Set(constants.HeaderXOmnifulCorrelationID, correlationID)
			c.Writer.Header().Set(constants.HeaderXOmnifulCorrelationID, correlationID)
		}

		requestID := c.GetHeader(constants.HeaderXOmnifulRequestID)

		if requestID == "" {
			requestID = NewRequestID()
		}

		// Setting Request ID in newRelic
		c.Set(constants.HeaderXOmnifulRequestID, requestID)

		c.Writer.Header().Set(constants.HeaderXOmnifulRequestID, requestID)

		c.Next()
	}
}

func NewRequestID() string {
	return uuid.New().String()
}

// GetRequestID returns the request ID from context
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(constants.HeaderXOmnifulRequestID).(string)
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

// GetRequestID returns the request ID from context
func GetRequestIDForPostgresqlLogging(ctx context.Context) string {
	requestID, ok := ctx.Value(constants.HeaderXOmnifulRequestID).(string)
	if !ok {
		return ""
	}

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, constants.HeaderXOmnifulRequestID, requestID)
}

func GetSQSMessageRequestID(ctx context.Context, attributes map[string]string) string {
	requestID, ok := attributes[constants.HeaderXOmnifulRequestID]
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

func SetSqsMessageRequestID(ctx context.Context, attributes map[string]string) context.Context {
	requestID := GetSQSMessageRequestID(ctx, attributes)
	// Setting newrelic request ID
	return context.WithValue(ctx, constants.HeaderXOmnifulRequestID, requestID)
}

func GetKafkaRequestID(ctx context.Context, headers map[string]string) string {
	requestID, ok := headers[constants.HeaderXOmnifulRequestID]
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

func SetKafkaRequestID(ctx context.Context, attributes map[string]string) context.Context {
	requestID := GetKafkaRequestID(ctx, attributes)

	// Setting newrelic request ID
	return context.WithValue(ctx, constants.HeaderXOmnifulRequestID, requestID)
}
