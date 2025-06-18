package db

import (
	"log"
	"os"

	"github.com/omniful/go_commons/db/sql/migration"
)


func RunMigrations() {
	dbURL := migration.BuildSQLDBURL(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	
	migrationPath := "file://migrations"

	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migration tool: %v", err)
	}

	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}
