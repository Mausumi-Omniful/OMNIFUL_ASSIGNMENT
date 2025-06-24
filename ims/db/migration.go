package db

import (
	"context"
	"log"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/migration"
)

func RunMigrations(ctx context.Context) {
	dbURL := migration.BuildSQLDBURL(
		config.GetString(ctx, "DB_HOST"),
		config.GetString(ctx, "DB_PORT"),
		config.GetString(ctx, "DB_NAME"),
		config.GetString(ctx, "DB_USER"),
		config.GetString(ctx, "DB_PASSWORD"),
	)

	migrationPath := "file://migrations"

	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	err = migrator.Up()
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}