package db

import (
	"errors"
	"os"
	"time"

	"github.com/omniful/go_commons/db/sql/postgres"
)

var DB *postgres.DbCluster

func InitPostgres() error {
	master := postgres.DBConfig{
		Host:                   os.Getenv("DB_HOST"),
		Port:                   os.Getenv("DB_PORT"),
		Username:               os.Getenv("DB_USER"),
		Password:               os.Getenv("DB_PASSWORD"),
		Dbname:                 os.Getenv("DB_NAME"),
		MaxOpenConnections:     10,
		MaxIdleConnections:     5,
		ConnMaxLifetime:        time.Minute * 5,
		DebugMode:              true,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	// Basic validation
	if master.Host == "" || master.Port == "" || master.Username == "" || master.Password == "" || master.Dbname == "" {
		return errors.New("missing DB config in .env")
	}

	var slaves []postgres.DBConfig // no slaves for now
	DB = postgres.InitializeDBInstance(master, &slaves)
	return nil
}

// GetDB allows external access
func GetDB() *postgres.DbCluster {
	return DB
}
