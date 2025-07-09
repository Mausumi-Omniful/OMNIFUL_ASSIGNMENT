package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKafkaProducer_PublishOrderCreated_Integration(t *testing.T) {

	brokers := []string{"localhost:9092"}
	topic := "order.created"

	producer, err := NewKafkaProducer(brokers, topic)
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	defer producer.Close()

	ctx := context.Background()
	event := OrderCreatedEvent{
		OrderID:   "order-123",
		SKU:       "sku-456",
		Location:  "loc-789",
		TenantID:  "tenant-001",
		SellerID:  "seller-002",
		Status:    "created",
		CreatedAt: "2024-01-01T12:00:00Z",
	}

	err = producer.PublishOrderCreated(ctx, event)
	assert.NoError(t, err)
}
