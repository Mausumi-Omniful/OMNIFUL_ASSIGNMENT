// main.go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"
	"github.com/Mausumi-Omniful/ims/routes"

	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("⚠️ Warning: Error loading .env file:", err)
	} else {
		fmt.Println("✅ .env loaded successfully")
	}

	// Initialize PostgreSQL connection
	if err := db.InitPostgres(); err != nil {
		fmt.Println("❌ Postgres connection failed:", err)
		panic("Postgres connection failed: " + err.Error())
	} else {
		fmt.Println("✅ Postgres connected successfully")
	}

	// AutoMigrate Inventory model
	if err := db.DB.GetMasterDB(context.Background()).AutoMigrate(
	&models.Inventory{},
	&models.SKU{},
	&models.Hub{},
); err != nil {
	panic("❌ AutoMigration failed: " + err.Error())
} else {
	fmt.Println("✅ Tables migrated: Inventory, SKU, Hub")
}

	// Initialize HTTP server
	server := http.InitializeServer(
		":8083",
		10*time.Second,
		10*time.Second,
		70*time.Second,
		false,
	)

	fmt.Println("🚀 IMS server is initializing on port 8083...")

	// Health check route
	// server.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })

	// Register inventory routes - pass *http.Server, NOT *gin.Engine
	routes.RegisterRoutes(server)

	// Start server
	if err := server.StartServer("ims-service"); err != nil {
		panic(err)
	}
}
