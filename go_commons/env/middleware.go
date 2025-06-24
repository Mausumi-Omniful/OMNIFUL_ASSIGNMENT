package env

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
)

// Middleware adds env in ctx
func Middleware(e interface{}) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set(constants.Env, e)
		c.Next()
	}
}
