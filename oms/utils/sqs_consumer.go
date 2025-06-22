package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"oms/database"
	"oms/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/omniful/go_commons/log"
)

type SQSConsumerImpl struct {
	client    *sqs.SQS
	queueURL  string
	queueName string
	handler   SQSMessageHandler
	isRunning bool
	stopChan  chan bool
}

type SQSMessageHandler interface {
	ProcessMessage(ctx context.Context, message *ConsumerMessage) error
}

type ConsumerMessage struct {
	RequestID string `json:"request_id"`
	Path      string `json:"path"`
	GroupID   string `json:"group_id"`
}

func NewSQSConsumer(queueName, endpoint, region string, handler SQSMessageHandler) (*SQSConsumerImpl, error) {
	// Set environment variables for LocalStack
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_DEFAULT_REGION", region)

	// Ensure queue name ends with .fifo for FIFO queue
	if !strings.HasSuffix(queueName, ".fifo") {
		queueName = queueName + ".fifo"
	}

	// Create AWS session for LocalStack
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"test", "test", "",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create SQS client
	sqsClient := sqs.New(sess)

	// Get queue URL
	queueURL := fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/%s", region, queueName)

	log.Infof("✅ SQS FIFO Consumer initialized for queue: %s", queueName)

	return &SQSConsumerImpl{
		client:    sqsClient,
		queueURL:  queueURL,
		queueName: queueName,
		handler:   handler,
		stopChan:  make(chan bool),
	}, nil
}

// Start begins consuming messages from the SQS queue
func (s *SQSConsumerImpl) Start(ctx context.Context) {
	if s.isRunning {
		log.Warn("Consumer is already running")
		return
	}

	s.isRunning = true
	log.Infof("🚀 Starting SQS Consumer for queue: %s", s.queueName)

	go s.consumeMessages(ctx)
}

// Stop gracefully stops the consumer
func (s *SQSConsumerImpl) Stop() error {
	if !s.isRunning {
		return nil
	}

	log.Infof("🛑 Stopping SQS Consumer for queue: %s", s.queueName)
	s.isRunning = false
	s.stopChan <- true
	return nil
}

// consumeMessages continuously polls for messages
func (s *SQSConsumerImpl) consumeMessages(ctx context.Context) {
	for s.isRunning {
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, stopping consumer")
			return
		case <-s.stopChan:
			log.Info("Stop signal received, stopping consumer")
			return
		default:
			s.pollMessages(ctx)
		}
	}
}

// pollMessages polls for messages and processes them
func (s *SQSConsumerImpl) pollMessages(ctx context.Context) {
	// Receive messages from SQS
	result, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: aws.Int64(1),  // Process one message at a time for FIFO
		WaitTimeSeconds:     aws.Int64(20), // Long polling
		VisibilityTimeout:   aws.Int64(30), // 30 seconds visibility timeout
	})

	if err != nil {
		log.WithError(err).Error("Failed to receive messages from SQS")
		time.Sleep(5 * time.Second) // Wait before retrying
		return
	}

	if len(result.Messages) == 0 {
		return // No messages available
	}

	// Process each message
	for _, message := range result.Messages {
		if err := s.processMessage(ctx, message); err != nil {
			log.WithError(err).Error("Failed to process message")
			// Don't delete the message, let it return to queue for retry
			continue
		}

		// Delete the message after successful processing
		if err := s.deleteMessage(message); err != nil {
			log.WithError(err).Error("Failed to delete message")
		}
	}
}

// processMessage processes a single SQS message
func (s *SQSConsumerImpl) processMessage(ctx context.Context, message *sqs.Message) error {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("❌ PANIC RECOVERED in processMessage: %v", r)
		}
	}()

	if message.Body == nil {
		return fmt.Errorf("message body is nil")
	}

	// Parse the JSON message
	var consumerMessage ConsumerMessage
	if err := json.Unmarshal([]byte(*message.Body), &consumerMessage); err != nil {
		return fmt.Errorf("failed to parse message JSON: %w", err)
	}

	log.Infof("📥 Processing message - RequestID: %s, Path: %s, GroupID: %s",
		consumerMessage.RequestID, consumerMessage.Path, consumerMessage.GroupID)

	// Process the message using the handler
	if err := s.handler.ProcessMessage(ctx, &consumerMessage); err != nil {
		return fmt.Errorf("handler failed to process message: %w", err)
	}

	log.Infof("✅ Successfully processed message - RequestID: %s", consumerMessage.RequestID)
	return nil
}

// deleteMessage deletes a message from the queue
func (s *SQSConsumerImpl) deleteMessage(message *sqs.Message) error {
	_, err := s.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}

// GetQueueName returns the queue name
func (s *SQSConsumerImpl) GetQueueName() string {
	return s.queueName
}

// DefaultMessageHandler provides a basic implementation for testing
type DefaultMessageHandler struct {
	s3Downloader  *S3DownloaderImpl
	csvParser     *CSVParser
	orderRepo     *database.OrderRepository
	imsClient     *IMSClient
	validator     *CSVRowValidator
	kafkaProducer *KafkaProducer
}

// NewDefaultMessageHandler creates a new default message handler with S3 and CSV capabilities
func NewDefaultMessageHandler(s3Endpoint, s3Region string, orderRepo *database.OrderRepository, imsClient *IMSClient, kafkaProducer *KafkaProducer, s3Uploader *S3UploaderImpl) (*DefaultMessageHandler, error) {
	// Create S3 downloader
	s3Downloader, err := NewS3Downloader(s3Endpoint, s3Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 downloader: %w", err)
	}

	// Create CSV parser
	csvParser := NewCSVParser(50) // batch size of 50

	// Create CSV row validator
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
			log.Errorf("PANIC RECOVERED in DefaultMessageHandler: %v", r)
		}
	}()

	if message == nil {
		return fmt.Errorf("message is nil")
	}
	if message.RequestID == "" {
		return fmt.Errorf("request_id is empty")
	}
	if message.Path == "" {
		return fmt.Errorf("path is empty")
	}

	log.Infof("Processing CSV file - RequestID: %s, S3 Path: %s", message.RequestID, message.Path)

	csvData, err := d.s3Downloader.DownloadFile(ctx, message.Path)
	if err != nil {
		return fmt.Errorf("failed to download CSV file from S3: %w", err)
	}

	parseResult, err := d.csvParser.ParseCSVFromBytes(ctx, csvData)
	if err != nil {
		return fmt.Errorf("failed to parse CSV data: %w", err)
	}

	log.Infof("CSV Processing Results: Total rows: %d, Valid rows: %d, Invalid rows: %d", parseResult.TotalRows, parseResult.ValidRows, parseResult.InvalidRows)

	if parseResult.ValidRows > 0 {
		createdOrders := 0
		failedOrders := 0
		validatedOrders := 0
		validationFailedOrders := 0

		for _, row := range parseResult.ValidData {
			validationResult := d.validator.ValidateCSVRow(ctx, row)
			if !validationResult.IsValid {
				validationFailedOrders++
				continue
			}
			validatedOrders++
			order := models.NewOrder(row.SKU, row.Location, row.TenantID, row.SellerID)
			if !order.IsValid() {
				failedOrders++
				continue
			}
			if err := d.orderRepo.SaveOrder(ctx, order); err != nil {
				failedOrders++
				continue
			}
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
			createdOrders++
		}
		log.Infof("Order creation summary: %d validated, %d saved to MongoDB, %d validation failed, %d failed", validatedOrders, createdOrders, validationFailedOrders, failedOrders)
	}

	if parseResult.InvalidRows > 0 {
		log.Infof("Found %d invalid rows - these will be skipped", parseResult.InvalidRows)
	}

	log.Infof("CSV processing completed for RequestID: %s", message.RequestID)
	return nil
}

// extractFilenameFromS3Path extracts filename from S3 path
func extractFilenameFromS3Path(s3Path string) string {
	// S3 path format: bucket/key/filename.csv
	parts := strings.Split(s3Path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown.csv"
}
