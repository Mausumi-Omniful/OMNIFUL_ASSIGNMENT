package routes

import (
	"github.com/mausumi-ghadei-omniful/ims/controllers"
	"github.com/omniful/go_commons/http"
)

func RegisterRoutes(server *http.Server) {
	// Inventory routes
	inv := server.Group("/inventory")
	inv.POST("/", controllers.CreateInventory)
	inv.GET("/", controllers.GetInventories)
	inv.PUT("/:id", controllers.UpdateInventory)
	inv.DELETE("/:id", controllers.DeleteInventory)
	inv.POST("/upsert", controllers.UpsertInventory)
	inv.POST("/reduce", controllers.ReduceInventory)

	// sku routes
	sku := server.Group("/sku")
	sku.POST("/", controllers.CreateSKU)
	sku.GET("/", controllers.GetSKUs)
	sku.PUT("/:id", controllers.UpdateSKU)
	sku.DELETE("/:id", controllers.DeleteSKU)

	// hub routes
	hub := server.Group("/hub")
	hub.POST("/", controllers.CreateHub)
	hub.GET("/", controllers.GetHubs)
	hub.PUT("/:id", controllers.UpdateHub)
	hub.DELETE("/:id", controllers.DeleteHub)
}
