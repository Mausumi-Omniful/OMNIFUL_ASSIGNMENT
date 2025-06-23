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

	
	migrationPath:="file://migrations"

	migrator,err:= migration.InitializeMigrate(migrationPath, dbURL)
	if err!= nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	err = migrator.Up()
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}