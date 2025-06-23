package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)


type OrderRepository struct {
	db         *Database
	collection *mongo.Collection
}



func NewOrderRepository(db *Database) *OrderRepository {
	collection := db.GetCollection("orders")
	return &OrderRepository{
		db:         db,
		collection: collection,
	}
}
