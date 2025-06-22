package db

import (
	"errors"
	"os"
	"github.com/omniful/go_commons/db/sql/postgres"
)

var DB *postgres.DbCluster

func InitPostgres() error{
	master := postgres.DBConfig{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname: os.Getenv("DB_NAME"),
	}

	
	if master.Host=="" || master.Port=="" || master.Username=="" || master.Password=="" || master.Dbname=="" {
		return errors.New("missing DB config")
	}

	var slaves []postgres.DBConfig
	DB = postgres.InitializeDBInstance(master, &slaves)
	return nil
}
