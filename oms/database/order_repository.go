package database

import (
	"context"
	"fmt"
	"time"

	"oms/models"

	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveOrder saves an order to MongoDB
func (r *OrderRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	// Set creation timestamp if not set
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}

	// Set update timestamp
	order.UpdatedAt = time.Now()

	// Convert to BSON document
	doc := bson.M{
		"order_id":   order.ID, // Use order_id field to store our custom ID
		"sku":        order.SKU,
		"location":   order.Location,
		"tenant_id":  order.TenantID,
		"seller_id":  order.SellerID,
		"status":     order.Status,
		"created_at": order.CreatedAt,
		"updated_at": order.UpdatedAt,
	}

	// Insert the document
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.WithError(err).Error("❌ Failed to save order to MongoDB")
		return err
	}

	// Log successful save with both MongoDB ObjectID and our custom Order ID
	log.Infof("✅ Order saved to MongoDB - ObjectID: %v, OrderID: %s", result.InsertedID, order.ID)
	return nil
}

// GetOrders retrieves orders from MongoDB with pagination
func (r *OrderRepository) GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Create options for pagination
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

	// Execute the query
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.WithError(err).Error("❌ Failed to query orders from MongoDB")
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		log.WithError(err).Error("❌ Failed to decode orders from MongoDB")
		return nil, err
	}

	log.Infof("✅ Retrieved %d orders from MongoDB", len(orders))
	return orders, nil
}

// GetOrdersByFilter retrieves orders with filtering and pagination
func (r *OrderRepository) GetOrdersByFilter(ctx context.Context, filters map[string]string, limit, offset int) ([]models.Order, error) {
	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Build filter query
	filter := bson.M{}

	// Add filters if provided
	if tenantID, exists := filters["tenant_id"]; exists && tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	if sellerID, exists := filters["seller_id"]; exists && sellerID != "" {
		filter["seller_id"] = sellerID
	}
	if status, exists := filters["status"]; exists && status != "" {
		filter["status"] = status
	}

	// Create options for pagination
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

	// Execute the query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.WithError(err).Error("❌ Failed to query filtered orders from MongoDB")
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		log.WithError(err).Error("❌ Failed to decode filtered orders from MongoDB")
		return nil, err
	}

	log.Infof("✅ Retrieved %d filtered orders from MongoDB", len(orders))
	return orders, nil
}

// GetOrderByID retrieves a specific order by its ID from MongoDB
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	log.Infof("🔍 Looking for order with ID: %s", orderID)

	// Create filter to find order by order_id
	filter := bson.M{"order_id": orderID}

	// Execute the query
	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			log.Warnf("⚠️ Order not found with ID: %s", orderID)
			return nil, fmt.Errorf("order not found with ID: %s", orderID)
		}
		log.WithError(err).Errorf("❌ Failed to query order from MongoDB - OrderID: %s", orderID)
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	// Convert BSON result to Order model with proper date handling
	order := &models.Order{
		ID:       result["order_id"].(string),
		SKU:      result["sku"].(string),
		Location: result["location"].(string),
		TenantID: result["tenant_id"].(string),
		SellerID: result["seller_id"].(string),
		Status:   models.OrderStatus(result["status"].(string)),
	}

	// Handle CreatedAt date conversion
	if createdAt, ok := result["created_at"].(primitive.DateTime); ok {
		order.CreatedAt = createdAt.Time()
	} else if createdAt, ok := result["created_at"].(time.Time); ok {
		order.CreatedAt = createdAt
	} else {
		log.Warnf("⚠️ Unexpected created_at type for order %s: %T", orderID, result["created_at"])
		order.CreatedAt = time.Now() // fallback
	}

	// Handle UpdatedAt date conversion
	if updatedAt, ok := result["updated_at"].(primitive.DateTime); ok {
		order.UpdatedAt = updatedAt.Time()
	} else if updatedAt, ok := result["updated_at"].(time.Time); ok {
		order.UpdatedAt = updatedAt
	} else {
		log.Warnf("⚠️ Unexpected updated_at type for order %s: %T", orderID, result["updated_at"])
		order.UpdatedAt = time.Now() // fallback
	}

	log.Infof("✅ Found order - OrderID: %s, Status: %s, SKU: %s, Location: %s",
		order.ID, order.Status, order.SKU, order.Location)
	return order, nil
}

// UpdateOrderStatus atomically updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus) error {
	log.Infof("🔄 Updating order status - OrderID: %s, NewStatus: %s", orderID, newStatus)

	// Create filter to find the specific order
	filter := bson.M{"order_id": orderID}

	// Create update with atomic operation
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}

	// Execute atomic update
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.WithError(err).Errorf("❌ Failed to update order status - OrderID: %s", orderID)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Check if order was found and updated
	if result.MatchedCount == 0 {
		log.Warnf("⚠️ Order not found for status update - OrderID: %s", orderID)
		return fmt.Errorf("order not found with ID: %s", orderID)
	}

	if result.ModifiedCount == 0 {
		log.Warnf("⚠️ Order status was not modified (already in target status) - OrderID: %s, Status: %s", orderID, newStatus)
		// This is not an error - the order might already be in the target status
	}

	log.Infof("✅ Order status updated successfully - OrderID: %s, NewStatus: %s, Modified: %d",
		orderID, newStatus, result.ModifiedCount)
	return nil
}
