package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db         *Database
	collection *mongo.Collection
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *Database) *OrderRepository {
	collection := db.GetCollection("orders")
	return &OrderRepository{
		db:         db,
		collection: collection,
	}
}
