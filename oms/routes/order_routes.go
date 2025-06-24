package routes

import (
	"oms/controllers"
	"oms/middleware"
	"oms/webhook"

	"github.com/omniful/go_commons/http"
)

// RegisterOrderRoutes
func RegisterOrderRoutes(server *http.Server, orderController *controllers.OrderController) {
	//middleware
	server.Use(middleware.LoggingMiddleware())


	// Order management routes
	orders := server.Group("/api/v1/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/upload", orderController.UploadCSV)
		orders.GET("/", orderController.ListOrders)
		orders.GET("/:orderID", orderController.GetOrderByID)
		orders.PUT("/:orderID/status", orderController.UpdateOrderStatus)
	}

	
	// Webhook events endpoint
	server.GET("/api/v1/webhook/events", webhook.GetWebhookEvents)
}
