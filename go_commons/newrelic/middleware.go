package newrelic

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()

			return
		}

		if clientService := c.GetHeader(constants.HeaderXClientService); len(clientService) > 0 {
			AddAttributeWithContext(
				c,
				Attribute{
					Name:  constants.HeaderXClientService,
					Value: clientService,
				},
			)
		}

		SetRequestID(c, env.GetRequestID(c))
		c.Next()
	}
}
