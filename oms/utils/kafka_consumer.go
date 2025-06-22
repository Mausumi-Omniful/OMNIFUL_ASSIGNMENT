package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"oms/models"

	"github.com/IBM/sarama"
	"github.com/omniful/go_commons/log"
)

// OrderRepositoryInterface defines the interface for order repository operations
type OrderRepositoryInterface interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error)
	GetOrdersByFilter(ctx context.Context, filters map[string]string, limit, offset int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus) error
}

// IMSClientInterface defines the interface for IMS client operations
type IMSClientInterface interface {
	GetSKUs() ([]SKU, error)
	GetHubs() ([]Hub, error)
	GetInventory() ([]Inventory, error)
	ValidateSKU(skuCode, tenantID, sellerID string) (bool, error)
	ValidateHub(hubName, tenantID, sellerID string) (bool, error)
	CheckInventoryAvailability(skuCode, location, tenantID, sellerID string) (bool, int, error)
}

// OrderFinalizationConsumer handles Kafka consumption for order finalization
type OrderFinalizationConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *OrderFinalizationHandler
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// OrderFinalizationHandler implements sarama.ConsumerGroupHandler for order finalization
type OrderFinalizationHandler struct {
	orderRepo OrderRepositoryInterface
	imsClient IMSClientInterface
}

// NewOrderFinalizationConsumer creates a new Kafka consumer for order finalization
func NewOrderFinalizationConsumer(brokers []string, topic string, orderRepo OrderRepositoryInterface, imsClient IMSClientInterface) (*OrderFinalizationConsumer, error) {
	// Create Sarama config
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "oms-order-finalization-group", config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Create order finalization handler
	handler := &OrderFinalizationHandler{
		orderRepo: orderRepo,
		imsClient: imsClient,
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	log.Infof("Order Finalization Kafka Consumer initialized for topic: %s", topic)

	return &OrderFinalizationConsumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// Start begins consuming messages from Kafka
func (c *OrderFinalizationConsumer) Start(ctx context.Context) {
	log.Infof("Starting Order Finalization Kafka Consumer...")

	// Add panic recovery for the consumer startup
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("❌ PANIC RECOVERED in Kafka Consumer Start: %v", r)
		}
	}()

	// Start consuming in a goroutine
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("❌ PANIC RECOVERED in Kafka Consumer goroutine: %v", r)
			}
		}()

		// Consume messages
		for {
			select {
			case <-c.ctx.Done():
				log.Infof("🛑 Kafka consumer context cancelled")
				return
			default:
				// Consume from the topic
				topics := []string{"order.created"}
				err := c.consumerGroup.Consume(c.ctx, topics, c.handler)
				if err != nil {
					log.WithError(err).Error("❌ Error in consumer group consume")
					// Continue trying to consume
				}
			}
		}
	}()
}

// Stop gracefully stops the consumer
func (c *OrderFinalizationConsumer) Stop() error {
	log.Infof("🛑 Stopping Order Finalization Kafka Consumer...")
	c.cancel()
	c.wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		log.WithError(err).Error("❌ Error closing consumer group")
		return err
	}
	return nil
}

// Setup implements sarama.ConsumerGroupHandler
func (h *OrderFinalizationHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Infof("Kafka consumer group session setup")
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler
func (h *OrderFinalizationHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Infof("Kafka consumer group session cleanup")
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler
func (h *OrderFinalizationHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Infof("Starting to consume messages from partition %d", claim.Partition())

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Infof("Message channel closed for partition %d", claim.Partition())
				return nil
			}

			// Process the message
			if err := h.processMessage(context.Background(), message); err != nil {
				log.WithError(err).Errorf("❌ Failed to process message from partition %d, offset %d",
					message.Partition, message.Offset)
				// Continue processing other messages
			} else {
				// Mark message as processed
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			log.Infof("Session context cancelled for partition %d", claim.Partition())
			return nil
		}
	}
}

// processMessage processes a single Kafka message
func (h *OrderFinalizationHandler) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("❌ PANIC RECOVERED in processMessage: %v", r)
		}
	}()

	log.Infof("Processing order.created event - Topic: %s, Key: %s", message.Topic, string(message.Key))

	// Parse the order.created event
	var orderEvent OrderCreatedEvent
	if err := json.Unmarshal(message.Value, &orderEvent); err != nil {
		log.WithError(err).Error("❌ Failed to parse order.created event")
		return fmt.Errorf("failed to parse order.created event: %w", err)
	}

	log.Infof("Processing order finalization for OrderID: %s, SKU: %s, Location: %s", orderEvent.OrderID, orderEvent.SKU, orderEvent.Location)

	// Step 1: Get the order from database
	order, err := h.orderRepo.GetOrderByID(ctx, orderEvent.OrderID)
	if err != nil {
		log.Infof("Failed to retrieve order for finalization - OrderID: %s", orderEvent.OrderID)
		return fmt.Errorf("failed to retrieve order: %w", err)
	}

	log.Infof("📋 Order retrieved - OrderID: %s, Current Status: %s", order.ID, order.Status)

	// Step 2: Check if order is in a valid state for finalization
	if order.Status != "on_hold" {
		log.Warnf("⚠️ Order is not in 'on_hold' status - OrderID: %s, Current Status: %s", order.ID, order.Status)
		return fmt.Errorf("order is not in 'on_hold' status: %s", order.Status)
	}

	// Step 3: Check inventory availability in IMS
	log.Infof("🔍 Checking inventory availability - SKU: %s, Location: %s, Tenant: %s, Seller: %s",
		order.SKU, order.Location, order.TenantID, order.SellerID)

	available, quantity, err := h.imsClient.CheckInventoryAvailability(order.SKU, order.Location, order.TenantID, order.SellerID)
	if err != nil {
		log.WithError(err).Errorf("❌ Failed to check inventory - OrderID: %s", order.ID)

		// Update order status to cancelled due to IMS error
		if updateErr := h.orderRepo.UpdateOrderStatus(ctx, order.ID, "cancelled"); updateErr != nil {
			log.WithError(updateErr).Errorf("❌ Failed to update order status to cancelled - OrderID: %s", order.ID)
		} else {
			log.Infof("✅ Order status updated to cancelled due to IMS error - OrderID: %s", order.ID)
		}

		return fmt.Errorf("failed to check inventory: %w", err)
	}

	// Step 4: Process based on inventory availability
	if available {
		// Inventory is available - finalize the order
		log.Infof("✅ Inventory available (Quantity: %d) - finalizing order - OrderID: %s", quantity, order.ID)

		if err := h.orderRepo.UpdateOrderStatus(ctx, order.ID, "new_order"); err != nil {
			log.WithError(err).Errorf("❌ Failed to update order status to new_order - OrderID: %s", order.ID)
			return fmt.Errorf("failed to finalize order: %w", err)
		}

		log.Infof("🎉 Order finalized successfully - OrderID: %s, Status: new_order", order.ID)
	} else {
		// Inventory not available - cancel the order
		log.Warnf("⚠️ Inventory not available (Quantity: %d) - cancelling order - OrderID: %s", quantity, order.ID)

		if err := h.orderRepo.UpdateOrderStatus(ctx, order.ID, "cancelled"); err != nil {
			log.WithError(err).Errorf("❌ Failed to update order status to cancelled - OrderID: %s", order.ID)
			return fmt.Errorf("failed to cancel order: %w", err)
		}

		log.Infof("❌ Order cancelled due to insufficient inventory - OrderID: %s, Status: cancelled", order.ID)
	}

	return nil
}
