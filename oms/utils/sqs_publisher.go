package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
)

type SQSPublisherImpl struct {
	client    *sqs.SQS
	queueURL  string
	queueName string
}

type SQSMessage struct {
	RequestID string `json:"request_id"`
	Path      string `json:"path"`
	GroupID   string `json:"group_id"`
}

func NewSQSPublisher(queueName, endpoint, region string) (*SQSPublisherImpl, error) {
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

	// Try to create the queue if it doesn't exist
	queueURL := fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/%s", region, queueName)

	// Check if queue exists by trying to get its attributes
	_, err = sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []*string{
			aws.String("All"),
		},
	})

	if err != nil {
		// Queue doesn't exist, create it
		log.Infof("Creating SQS FIFO queue: %s", queueName)

		createResult, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
			Attributes: map[string]*string{
				"FifoQueue":                 aws.String("true"),
				"ContentBasedDeduplication": aws.String("true"),
				"VisibilityTimeout":         aws.String("30"),
				"MessageRetentionPeriod":    aws.String("345600"), // 4 days
			},
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create SQS FIFO queue: %w", err)
		}

		queueURL = *createResult.QueueUrl
		log.Infof("SQS FIFO queue created: %s", queueName)
	} else {
		log.Infof("SQS FIFO queue already exists: %s", queueName)
	}

	log.Infof("SQS FIFO Publisher initialized for queue: %s", queueName)

	return &SQSPublisherImpl{
		client:    sqsClient,
		queueURL:  queueURL,
		queueName: queueName,
	}, nil
}

func (s *SQSPublisherImpl) PublishMessage(ctx context.Context, messageBody string) error {
	// Create SQS message
	message := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(messageBody),
	}

	// Publish message to SQS
	_, err := s.client.SendMessage(message)
	if err != nil {
		log.WithError(err).Error("❌ Failed to publish message to SQS")
		return fmt.Errorf("failed to publish message to SQS: %w", err)
	}

	log.Infof("Message published to SQS FIFO queue '%s': %s", s.queueName, messageBody)
	return nil
}

func (s *SQSPublisherImpl) PublishS3Path(ctx context.Context, s3Path string) error {
	// Generate unique request ID
	requestID := uuid.New().String()

	// Create structured message with group ID
	sqsMessage := SQSMessage{
		RequestID: requestID,
		Path:      s3Path,
		GroupID:   "csv-processing", // You can customize this based on tenant_id or other criteria
	}

	// Convert to JSON
	jsonMessage, err := json.Marshal(sqsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal SQS message: %w", err)
	}

	log.Infof("Publishing structured message to SQS FIFO - RequestID: %s, Path: %s, GroupID: %s", requestID, s3Path, sqsMessage.GroupID)

	// Create SQS message with group ID for FIFO
	message := &sqs.SendMessageInput{
		QueueUrl:               aws.String(s.queueURL),
		MessageBody:            aws.String(string(jsonMessage)),
		MessageGroupId:         aws.String(sqsMessage.GroupID),
		MessageDeduplicationId: aws.String(requestID), // Use requestID for deduplication
	}

	// Publish message to SQS FIFO
	_, err = s.client.SendMessage(message)
	if err != nil {
		log.WithError(err).Error("❌ Failed to publish message to SQS FIFO")
		return fmt.Errorf("failed to publish message to SQS FIFO: %w", err)
	}

	log.Infof("Message published to SQS FIFO queue '%s' with GroupID: %s", s.queueName, sqsMessage.GroupID)
	return nil
}

func (s *SQSPublisherImpl) GetQueueName() string {
	return s.queueName
}
