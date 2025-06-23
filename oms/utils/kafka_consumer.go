package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"oms/models"

	"github.com/IBM/sarama"
)

type OrderRepositoryInterface interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error)
	GetOrdersByFilter(ctx context.Context, filters map[string]string, limit, offset int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus) error
}

type IMSClientInterface interface {
	GetSKUs() ([]SKU, error)
	GetHubs() ([]Hub, error)
	GetInventory() ([]Inventory, error)
	ValidateSKU(skuCode, tenantID, sellerID string) (bool, error)
	ValidateHub(hubName, tenantID, sellerID string) (bool, error)
	CheckInventoryAvailability(skuCode, location, tenantID, sellerID string) (bool, int, error)
}

type OrderFinalizationConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *OrderFinalizationHandler
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

type OrderFinalizationHandler struct {
	orderRepo OrderRepositoryInterface
	imsClient IMSClientInterface
}

func NewOrderFinalizationConsumer(brokers []string, topic string, orderRepo OrderRepositoryInterface, imsClient IMSClientInterface) (*OrderFinalizationConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "oms-order-finalization-group", config)
	if err != nil {
		return nil, fmt.Errorf("new group error: %w", err)
	}

	handler := &OrderFinalizationHandler{orderRepo: orderRepo, imsClient: imsClient}
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Printf("Consumer ready for topic: %s\n", topic)

	return &OrderFinalizationConsumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (c *OrderFinalizationConsumer) Start(ctx context.Context) {
	fmt.Println("Consumer starting...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in Start: %v\n", r)
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in goroutine: %v\n", r)
			}
		}()

		for {
			select {
			case <-c.ctx.Done():
				fmt.Println("Context cancelled")
				return
			default:
				topics := []string{"order.created"}
				err := c.consumerGroup.Consume(c.ctx, topics, c.handler)
				if err != nil {
					fmt.Printf("Consume error: %v\n", err)
				}
			}
		}
	}()
}

func (c *OrderFinalizationConsumer) Stop() error {
	fmt.Println("Stopping consumer...")
	c.cancel()
	c.wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		fmt.Printf("Close error: %v\n", err)
		return err
	}
	return nil
}

func (h *OrderFinalizationHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("Setup done")
	return nil
}

func (h *OrderFinalizationHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Cleanup done")
	return nil
}



// consumeclaim
func (h *OrderFinalizationHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("Consuming from partition %d\n", claim.Partition())

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				fmt.Println("Channel closed")
				return nil
			}
			err := h.processMessage(context.Background(), message)
			if err != nil {
				fmt.Printf("Process error: %v\n", err)
			} else {
				session.MarkMessage(message, "")
			}
		case <-session.Context().Done():
			fmt.Println("Session done")
			return nil
		}
	}
}






// processmessage
func (h *OrderFinalizationHandler) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in processMessage: %v\n", r)
		}
	}()

	fmt.Printf("Received event: %s\n", string(message.Key))

	var orderEvent OrderCreatedEvent
	if err := json.Unmarshal(message.Value, &orderEvent); err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return fmt.Errorf("parse failed: %w", err)
	}

	fmt.Printf("Processing OrderID: %s\n", orderEvent.OrderID)

	order, err := h.orderRepo.GetOrderByID(ctx, orderEvent.OrderID)
	if err != nil {
		fmt.Println("Order fetch failed")
		return fmt.Errorf("fetch failed: %w", err)
	}

	fmt.Printf("Status: %s\n", order.Status)

	if order.Status != "on_hold" {
		fmt.Println("Invalid status")
		return fmt.Errorf("invalid status: %s", order.Status)
	}

	available, quantity, err := h.imsClient.CheckInventoryAvailability(order.SKU, order.Location, order.TenantID, order.SellerID)
	if err != nil {
		fmt.Println("Inventory check failed")
		_ = h.orderRepo.UpdateOrderStatus(ctx, order.ID, "cancelled")
		fmt.Println("Order cancelled due to error")
		return fmt.Errorf("inventory error: %w", err)
	}

	if available {
		fmt.Printf("Stock OK (%d), finalizing\n", quantity)
		if err := h.orderRepo.UpdateOrderStatus(ctx, order.ID, "new_order"); err != nil {
			fmt.Println("Status update failed")
			return fmt.Errorf("finalize error: %w", err)
		}
		fmt.Println("Order finalized")
	} else {
		fmt.Println("No stock, cancelling")
		if err := h.orderRepo.UpdateOrderStatus(ctx, order.ID, "cancelled"); err != nil {
			fmt.Println("Cancel failed")
			return fmt.Errorf("cancel error: %w", err)
		}
		fmt.Println("Order cancelled")
	}

	return nil
}
