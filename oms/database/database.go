package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalDB *Database

func SetGlobalDatabase(db *Database) {
	globalDB = db
}

func GetGlobalDatabase() *Database {
	return globalDB
}

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

// MongoDB connection
func NewDatabase(ctx context.Context, uri, dbName string) (*Database, error) {
	fmt.Println("Connecting to MongoDB:", uri)

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("MongoDB connected successfully to database:", dbName)

	return &Database{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (db *Database) Close(ctx context.Context) error {
	if db.client != nil {
		fmt.Println("Closing MongoDB connection...")
		if err := db.client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to close MongoDB connection: %w", err)
		}
		fmt.Println("MongoDB connection closed")
	}
	return nil
}
