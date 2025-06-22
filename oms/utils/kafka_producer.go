package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

// OrderCreatedEvent represents the event structure for order.created
type OrderCreatedEvent struct {
	OrderID   string `json:"order_id"`
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	TenantID  string `json:"tenant_id"`
	SellerID  string `json:"seller_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// KafkaProducer handles Kafka event publishing
type KafkaProducer struct {
	producer *kafka.ProducerClient
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	producer := kafka.NewProducer(
		kafka.WithBrokers(brokers),
		kafka.WithClientID("oms-service"),
		kafka.WithKafkaVersion("2.8.1"),
	)

	log.Infof("Kafka producer initialized for topic: %s", topic)

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishOrderCreated publishes an order.created event
func (k *KafkaProducer) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
	// Convert event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal order created event: %w", err)
	}

	// Create pubsub message
	message := &pubsub.Message{
		Topic: k.topic,
		Value: eventData,
		Key:   event.OrderID, // Use order ID as message key for partitioning
	}

	// Publish message
	err = k.producer.Publish(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to publish order created event: %w", err)
	}

	log.Infof("Published order.created event for order: %s", event.OrderID)
	return nil
}

// Close closes the Kafka producer
func (k *KafkaProducer) Close() {
	if k.producer != nil {
		k.producer.Close()
		log.Infof("Kafka producer closed")
	}
}
