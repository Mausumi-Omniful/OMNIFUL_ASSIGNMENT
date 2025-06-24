package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"oms/models"

	"github.com/omniful/go_commons/httpclient"
	"github.com/omniful/go_commons/httpclient/request"
)


type WebhookPayload struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	SKU       string `json:"sku"`
	Location  string `json:"location"`
	Timestamp string `json:"timestamp"`
}


func SendWebhook(ctx context.Context, order *models.Order, webhookURL string) error {
	if webhookURL == "" {
		return nil
	}

	payload := WebhookPayload{
		OrderID:   order.ID,
		Status:    string(order.Status),
		SKU:       order.SKU,
		Location:  order.Location,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}


	headers := url.Values{}
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "OMS-Webhook/1.0")


	req, err := request.NewBuilder().
		SetMethod("POST").
		SetUri(webhookURL).
		SetHeaders(headers).
		SetBody(jsonData).
		Build()
	if err != nil {
		return fmt.Errorf("failed to build webhook request: %w", err)
	}
	client := httpclient.New("", httpclient.WithTimeout(10*time.Second))
	resp, err := client.Post(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("webhook failed with status: %d", resp.StatusCode())
	}

	fmt.Printf("Webhook sent successfully - OrderID: %s, Status: %s\n", order.ID, order.Status)
	return nil
}
