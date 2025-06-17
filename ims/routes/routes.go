package routes

import (
	"github.com/Mausumi-Omniful/ims/controllers"
	"github.com/omniful/go_commons/http"
)

func RegisterRoutes(server *http.Server) {

	server.POST("/inventories", controllers.CreateInventory)
	server.GET("/inventories", controllers.GetInventories)
    server.GET("/inventories/:id", controllers.GetInventoryByID)

}
