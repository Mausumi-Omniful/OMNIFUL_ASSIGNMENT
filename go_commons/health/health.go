package health

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/shutdown"
	"net/http"
	"sync/atomic"
)

type ginHealthHandler struct {
	healthy *uint32
}

func (gh *ginHealthHandler) Close() error {
	return nil
}

const healthyStatusCode = uint32(0)
const unhealthyStatusCode = uint32(1)

// HealthcheckHandler returns a gin.HandlerFunc that can be passed onto
// the gin router, once registered health check handler will respond to the health-checks
func HealthcheckHandler() gin.HandlerFunc {
	healthy := healthyStatusCode
	gh := &ginHealthHandler{
		healthy: &healthy,
	}

	// Register for drain callback
	shutdown.RegisterDrainCallback(constants.HealthCheck, gh)

	return func(c *gin.Context) {
		if atomic.LoadUint32(gh.healthy) == healthyStatusCode {
			c.String(http.StatusOK, constants.Healthy)
		} else {
			c.String(http.StatusGone, constants.Unhealthy)
		}
	}
}
