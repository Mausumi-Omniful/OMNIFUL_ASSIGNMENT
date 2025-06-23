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
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_DEFAULT_REGION", region)

	if !strings.HasSuffix(queueName, ".fifo") {
		queueName += ".fifo"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return nil, fmt.Errorf("AWS session error: %w", err)
	}

	sqsClient := sqs.New(sess)
	queueURL := fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/%s", region, queueName)

	_, err = sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []*string{
			aws.String("All"),
		},
	})

	if err != nil {
		fmt.Printf("Creating queue: %s\n", queueName)
		createResult, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String(queueName),
			Attributes: map[string]*string{
				"FifoQueue":                 aws.String("true"),
				"ContentBasedDeduplication": aws.String("true"),
				"VisibilityTimeout":         aws.String("30"),
				"MessageRetentionPeriod":    aws.String("345600"),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("queue creation error: %w", err)
		}
		queueURL = *createResult.QueueUrl
		fmt.Println("Queue created")
	} else {
		fmt.Println("Queue exists")
	}

	fmt.Println("Publisher ready")

	return &SQSPublisherImpl{
		client:    sqsClient,
		queueURL:  queueURL,
		queueName: queueName,
	}, nil
}






func (s *SQSPublisherImpl) PublishS3Path(ctx context.Context, s3Path string) error {
	requestID := uuid.New().String()

	msg := SQSMessage{
		RequestID: requestID,
		Path:      s3Path,
		GroupID:   "csv-processing",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	fmt.Printf("Sending: %s\n", requestID)

	message := &sqs.SendMessageInput{
		QueueUrl:               aws.String(s.queueURL),
		MessageBody:            aws.String(string(data)),
		MessageGroupId:         aws.String(msg.GroupID),
		MessageDeduplicationId: aws.String(requestID),
	}

	_, err = s.client.SendMessage(message)
	if err != nil {
		fmt.Println("FIFO publish failed")
		return fmt.Errorf("FIFO publish error: %w", err)
	}

	fmt.Println("FIFO message sent")
	return nil
}




func (s *SQSPublisherImpl) GetQueueName() string {
	return s.queueName
}
