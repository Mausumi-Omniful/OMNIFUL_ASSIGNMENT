package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/redisclient"
	"github.com/Mausumi-Omniful/ims/routes"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
)

func main() {
	// Set CONFIG_SOURCE=local if not already set
	if os.Getenv("CONFIG_SOURCE") == "" {
		os.Setenv("CONFIG_SOURCE", "local")
	}

	// Initialize go_commons/config
	err := config.Init(15 * time.Second)
	if err != nil {
		fmt.Println("Failed to initialize config:", err)
		os.Exit(1)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		fmt.Println("Failed to get config context:", err)
		os.Exit(1)
	}

	// Initialize PostgreSQL connection
	err = db.InitPostgres(ctx)
	if err != nil {
		fmt.Println("Postgres connection failed:", err)
	} else {
		fmt.Println("Postgres connected successfully")
	}

	// migrations
	db.RunMigrations(ctx)

	// Initialize Redis client
	if err := redisclient.InitRedis(ctx); err != nil {
		fmt.Println("Redis initialization failed:", err)
	} else {
		fmt.Println("Redis connected successfully")
	}

	defer func() {
		if err := redisclient.Close(); err != nil {
			fmt.Println("Error closing Redis client:", err)
		} else {
			fmt.Println("Redis client closed successfully")
		}
	}()

	// Initialize HTTP server
	server := http.InitializeServer(
		":8084",
		10*time.Second,
		10*time.Second,
		70*time.Second,
		false,
	)

	fmt.Println("IMS server is running on port 8084...")

	// Register routes
	routes.RegisterRoutes(server)

	// Start server
	if err := server.StartServer("ims-service"); err != nil {
		panic(err)
	}
}