package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMonitoringRoutes(server *http.Server) {
	server.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus Monitoring
	// TODO add more monitoring routes like memory, cpu profiling
}
