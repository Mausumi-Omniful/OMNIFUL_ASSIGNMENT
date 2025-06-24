package db

import (
	"context"
	"errors"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/postgres"
)

var DB *postgres.DbCluster

func InitPostgres(ctx context.Context) error {
	master := postgres.DBConfig{
		Host:     config.GetString(ctx, "DB_HOST"),
		Port:     config.GetString(ctx, "DB_PORT"),
		Username: config.GetString(ctx, "DB_USER"),
		Password: config.GetString(ctx, "DB_PASSWORD"),
		Dbname:   config.GetString(ctx, "DB_NAME"),
	}

	if master.Host == "" || master.Port == "" || master.Username == "" || master.Password == "" || master.Dbname == "" {
		return errors.New("missing DB config")
	}

	var slaves []postgres.DBConfig
	DB = postgres.InitializeDBInstance(master, &slaves)
	return nil
}