package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)


type OrderRepository struct {
	db         *Database
	collection *mongo.Collection
}

func (db *Database) GetCollection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

func NewOrderRepository(db *Database) *OrderRepository {
	collection := db.GetCollection("orders")
	return &OrderRepository{
		db:         db,
		collection: collection,
	}
}
