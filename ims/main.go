package main

import (
	"fmt"
	"time"
    
	"github.com/Mausumi-Omniful/ims/db"
	"github.com/Mausumi-Omniful/ims/redisclient"
	"github.com/Mausumi-Omniful/ims/routes"

	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
)

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err!= nil {
		fmt.Println(err)
	} else {
		fmt.Println(".env loaded")
	}
	


	// Initialize PostgreSQL connection
	err= db.InitPostgres()
	if err != nil {
		fmt.Println("Postgres connection failed:", err)
	} else {
		fmt.Println("Postgres connected successfully")
	}



	// migrations
     db.RunMigrations()

  

	// Initialize Redis client
	if err := redisclient.InitRedis(); err != nil {
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

	fmt.Println("IMS server is running on port 8083...")

	
	// Register routes
	routes.RegisterRoutes(server)


	// Start server
	if err := server.StartServer("ims-service"); err != nil {
		panic(err)
	}
}
