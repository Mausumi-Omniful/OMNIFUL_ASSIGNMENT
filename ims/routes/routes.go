package routes

import (
	"github.com/Mausumi-Omniful/ims/controllers"
	"github.com/omniful/go_commons/http"
)

func RegisterRoutes(server *http.Server) {
	// Inventory routes
	inv := server.Group("/inventory")
    inv.GET("/", controllers.GetInventories)
	inv.GET("/:id", controllers.GetInventoryByID)
	inv.PUT("/:id", controllers.UpdateInventory)
	inv.POST("/upsert", controllers.UpsertInventory)
	inv.DELETE("/:id", controllers.DeleteInventory)


	
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
	hub.GET("/:id", controllers.GetHubByID)
	hub.PUT("/:id", controllers.UpdateHub)
	hub.DELETE("/:id", controllers.DeleteHub)
}
