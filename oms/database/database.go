package database

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents the MongoDB connection and database instance
type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewDatabase creates a new MongoDB connection
func NewDatabase(ctx context.Context, uri, dbName string) (*Database, error) {
	log.Infof("🔄 Connecting to MongoDB: %s", uri)

	// Set connection options
	clientOptions := options.Client().ApplyURI(uri)

	// Set timeout for connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database instance
	database := client.Database(dbName)

	log.Infof("✅ MongoDB connected successfully to database: %s", dbName)

	return &Database{
		client:   client,
		database: database,
	}, nil
}

// GetCollection returns a MongoDB collection by name
func (db *Database) GetCollection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

// Close closes the MongoDB connection
func (db *Database) Close(ctx context.Context) error {
	if db.client != nil {
		log.Infof("🔄 Closing MongoDB connection...")
		err := db.client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("failed to close MongoDB connection: %w", err)
		}
		log.Infof("✅ MongoDB connection closed")
	}
	return nil
}

// GetClient returns the MongoDB client (for advanced operations)
func (db *Database) GetClient() *mongo.Client {
	return db.client
}

// GetDatabase returns the MongoDB database instance
func (db *Database) GetDatabase() *mongo.Database {
	return db.database
}
