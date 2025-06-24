package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_DEFAULT_REGION", region)

	if !strings.HasSuffix(queueName, ".fifo") {
		queueName = queueName + ".fifo"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	sqsClient := sqs.New(sess)
	queueURL := fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/%s", region, queueName)

	fmt.Printf("SQS FIFO Consumer initialized for queue: %s\n", queueName)

	return &SQSConsumerImpl{
		client:    sqsClient,
		queueURL:  queueURL,
		queueName: queueName,
		handler:   handler,
		stopChan:  make(chan bool),
	}, nil
}




func (s *SQSConsumerImpl) Start(ctx context.Context) {
	if s.isRunning {
		fmt.Println("Consumer is already running")
		return
	}

	s.isRunning = true
	fmt.Printf("Starting SQS Consumer for queue: %s\n", s.queueName)

	go s.consumeMessages(ctx)
}



func (s *SQSConsumerImpl) consumeMessages(ctx context.Context) {
	for s.isRunning {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping consumer")
			return
		case <-s.stopChan:
			fmt.Println("Stop signal received, stopping consumer")
			return
		default:
			s.pollMessages(ctx)
		}
	}
}




func (s *SQSConsumerImpl) pollMessages(ctx context.Context) {
	result, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(20),
		VisibilityTimeout:   aws.Int64(30),
	})

	if err != nil {
		fmt.Printf("Error receiving messages: %v\n", err)
		time.Sleep(5 * time.Second)
		return
	}

	if len(result.Messages) == 0 {
		return
	}

	for _, message := range result.Messages {
		if err := s.processMessage(ctx, message); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
			continue
		}
		if err := s.deleteMessage(message); err != nil {
			fmt.Printf("Error deleting message: %v\n", err)
		}
	}
}






func (s *SQSConsumerImpl) processMessage(ctx context.Context, message *sqs.Message) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered in processMessage: %v\n", r)
		}
	}()

	if message.Body == nil {
		return fmt.Errorf("message body is nil")
	}

	var consumerMessage ConsumerMessage
	if err := json.Unmarshal([]byte(*message.Body), &consumerMessage); err != nil {
		return fmt.Errorf("failed to parse message JSON: %w", err)
	}

	fmt.Printf("Processing message - RequestID: %s\n", consumerMessage.RequestID)

	if err := s.handler.ProcessMessage(ctx, &consumerMessage); err != nil {
		return fmt.Errorf("handler failed to process message: %w", err)
	}

	fmt.Printf("Message processed - RequestID: %s\n", consumerMessage.RequestID)
	return nil
}

func (s *SQSConsumerImpl) deleteMessage(message *sqs.Message) error {
	_, err := s.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}