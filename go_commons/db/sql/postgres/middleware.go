package postgres

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
)

// Middleware adds db consistency in request
func Middleware() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set(constants.Consistency, &Consistency{constants.EventualConsistency})
		c.Next()
	}
}

func SetStrongConsistency() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set(constants.Consistency, &Consistency{constants.StrongConsistency})
		c.Next()
	}
}

func SetSlaveDBPreference() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set(constants.DBPreference, constants.SlaveDB)
		c.Next()
	}
}
