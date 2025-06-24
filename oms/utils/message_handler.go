package utils

import (
	"context"
	"fmt"

	"oms/database"
	"oms/models"
	"oms/webhook"
)

type DefaultMessageHandler struct {
	s3Downloader  *S3DownloaderImpl
	csvParser     *CSVParser
	orderRepo     *database.OrderRepository
	imsClient     *IMSClient
	validator     *CSVRowValidator
	kafkaProducer *KafkaProducer
}

func NewDefaultMessageHandler(s3Endpoint, s3Region string, orderRepo *database.OrderRepository, imsClient *IMSClient, kafkaProducer *KafkaProducer, s3Uploader *S3UploaderImpl) (*DefaultMessageHandler, error) {
	s3Downloader, err := NewS3Downloader(s3Endpoint, s3Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 downloader: %w", err)
	}

	csvParser := NewCSVParser(50)
	validator := NewCSVRowValidator(imsClient)

	return &DefaultMessageHandler{
		s3Downloader:  s3Downloader,
		csvParser:     csvParser,
		orderRepo:     orderRepo,
		imsClient:     imsClient,
		validator:     validator,
		kafkaProducer: kafkaProducer,
	}, nil
}

func (d *DefaultMessageHandler) ProcessMessage(ctx context.Context, message *ConsumerMessage) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered in DefaultMessageHandler: %v\n", r)
		}
	}()

	if message == nil || message.RequestID == "" || message.Path == "" {
		return fmt.Errorf("invalid message content")
	}

	fmt.Printf("Downloading and processing CSV - RequestID: %s, Path: %s\n", message.RequestID, message.Path)

	csvData, err := d.s3Downloader.DownloadFile(ctx, message.Path)
	if err != nil {
		return fmt.Errorf("failed to download CSV: %w", err)
	}

	parseResult, err := d.csvParser.ParseCSVFromBytes(ctx, csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	for _, row := range parseResult.ValidData {
		if !d.validator.ValidateCSVRow(ctx, row).IsValid {
			continue
		}

		order := models.NewOrder(row.SKU, row.Location, row.TenantID, row.SellerID)
		if !order.IsValid() {
			continue
		}

		if err := d.orderRepo.SaveOrder(ctx, order); err != nil {
			continue
		}
		// Log webhook event for order creation
		_ = webhook.LogWebhookEvent(ctx, "order.created", order)

		if d.kafkaProducer != nil {
			event := OrderCreatedEvent{
				OrderID:   order.ID,
				SKU:       order.SKU,
				Location:  order.Location,
				TenantID:  order.TenantID,
				SellerID:  order.SellerID,
				Status:    string(order.Status),
				CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			_ = d.kafkaProducer.PublishOrderCreated(ctx, event)
		}
	}

	fmt.Printf("CSV processing completed - RequestID: %s\n", message.RequestID)
	return nil
}
