package routes

import (
	"oms/controllers"
	"oms/middleware"

	"github.com/omniful/go_commons/http"
)

// RegisterOrderRoutes registers all order-related routes using go_commons patterns
func RegisterOrderRoutes(server *http.Server, orderController *controllers.OrderController) {
	// Apply global middleware
	server.Use(middleware.LoggingMiddleware())

	// Public endpoints (no auth required)
	server.GET("/test", orderController.TestEndpoint)
	server.GET("/health", orderController.HealthCheck)

	// Order management routes (with auth required)
	orders := server.Group("/api/v1/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/upload", orderController.UploadCSV)
		orders.GET("/", orderController.ListOrders)
		orders.GET("/:orderID", orderController.GetOrderByID)
		orders.PUT("/:orderID/status", orderController.UpdateOrderStatus)
	}
}
