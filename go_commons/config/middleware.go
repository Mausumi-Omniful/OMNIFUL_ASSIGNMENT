package config

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
)

// Middleware ensures that config will remain same for a single request in his whole journey
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tempApp := getApplication()
		if tempApp == nil {
			c.Abort()
		}

		conf := tempApp.observer.GetConfig()
		c.Set(constants.Config, conf)
		c.Next()
	}
}
