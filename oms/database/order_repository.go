package database

import (
	"context"
	"fmt"
	"time"

	"oms/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)





//saves order to MongoDB
func (r *OrderRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	order.UpdatedAt = time.Now()

	doc := bson.M{
		"order_id":   order.ID, 
		"sku":        order.SKU,
		"location":   order.Location,
		"tenant_id":  order.TenantID,
		"seller_id":  order.SellerID,
		"status":     order.Status,
		"created_at": order.CreatedAt,
		"updated_at": order.UpdatedAt,
	}

	
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		fmt.Println("ERROR: Failed to save order to MongoDB:", err)
		return err
	}
	fmt.Printf("Order saved to MongoDB - ObjectID: %v, OrderID: %s\n", result.InsertedID, order.ID)
	return nil
}







// GetOrders from mongodb
func (r *OrderRepository) GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "created_at", Value: -1}}) 
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		fmt.Println("ERROR: Failed to query orders from MongoDB:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		fmt.Println("ERROR: Failed to decode orders from MongoDB:", err)
		return nil, err
	}

	fmt.Printf("Retrieved %d orders from MongoDB\n", len(orders))
	return orders, nil
}






// GetOrdersByFilter
func (r *OrderRepository) GetOrdersByFilter(ctx context.Context, filters map[string]string, limit, offset int) ([]models.Order, error) {

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	filter := bson.M{}

	if tenantID, exists := filters["tenant_id"]; exists && tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	if sellerID, exists := filters["seller_id"]; exists && sellerID != "" {
		filter["seller_id"] = sellerID
	}
	if status, exists := filters["status"]; exists && status != "" {
		filter["status"] = status
	}

	
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "created_at", Value: -1}}) 
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		fmt.Println("ERROR: Failed to query filtered orders from MongoDB:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		fmt.Println("ERROR: Failed to decode filtered orders from MongoDB:", err)
		return nil, err
	}

	fmt.Printf("Retrieved %d filtered orders from MongoDB\n", len(orders))
	return orders, nil
}






// GetOrderByID
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	fmt.Printf("Looking for order with ID: %s\n", orderID)
	filter := bson.M{"order_id": orderID}
	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Printf("Order not found with ID: %s\n", orderID)
			return nil, fmt.Errorf("order not found with ID: %s", orderID)
		}
		fmt.Printf("ERROR: Failed to query order from MongoDB - OrderID: %s: %v\n", orderID, err)
		return nil, fmt.Errorf("failed to query order: %w", err)
	}
	order := &models.Order{
		ID:       result["order_id"].(string),
		SKU:      result["sku"].(string),
		Location: result["location"].(string),
		TenantID: result["tenant_id"].(string),
		SellerID: result["seller_id"].(string),
		Status:   models.OrderStatus(result["status"].(string)),
	}
	if createdAt, ok := result["created_at"].(primitive.DateTime); ok {
		order.CreatedAt = createdAt.Time()
	} else if createdAt, ok := result["created_at"].(time.Time); ok {
		order.CreatedAt = createdAt
	} else {
		order.CreatedAt = time.Now() // fallback
	}
	if updatedAt, ok := result["updated_at"].(primitive.DateTime); ok {
		order.UpdatedAt = updatedAt.Time()
	} else if updatedAt, ok := result["updated_at"].(time.Time); ok {
		order.UpdatedAt = updatedAt
	} else {
		order.UpdatedAt = time.Now() // fallback
	}

	fmt.Printf("Found order - OrderID: %s, Status: %s, SKU: %s, Location: %s\n",
		order.ID, order.Status, order.SKU, order.Location)
	return order, nil
}










// UpdateOrderStatus
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus) error {
	fmt.Printf("Updating order status - OrderID: %s, NewStatus: %s\n", orderID, newStatus)
	filter := bson.M{"order_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("ERROR: Failed to update order status - OrderID: %s: %v\n", orderID, err)
		return fmt.Errorf("failed to update order status: %w", err)
	}
	if result.MatchedCount == 0 {
		fmt.Printf("Order not found for status update - OrderID: %s\n", orderID)
		return fmt.Errorf("order not found with ID: %s", orderID)
	}

	if result.ModifiedCount == 0 {
		fmt.Printf("Order status was not modified (already in target status) - OrderID: %s, Status: %s\n", orderID, newStatus)
	}

	fmt.Printf("Order status updated successfully - OrderID: %s, NewStatus: %s, Modified: %d\n",
		orderID, newStatus, result.ModifiedCount)
	return nil
}
