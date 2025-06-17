package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"github.com/omniful/go_commons/http"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load("../.env") // Adjust path if needed
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

	// AutoMigrate Inventory model to create table in DB
	if err := db.DB.GetMasterDB(context.Background()).AutoMigrate(&models.Inventory{}); err != nil {
	panic("❌ AutoMigration failed: " + err.Error())
} else {
		fmt.Println("✅ Inventory table migrated successfully")
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

	// Add a basic health check route
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Start the HTTP server
	if err := server.StartServer("ims-service"); err != nil {
		panic(err)
	}
}
