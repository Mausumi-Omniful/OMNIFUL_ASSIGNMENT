package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

type QueueType string

const (
	QueueStandard QueueType = "standard"
	QueueFifo     QueueType = "fifo"
)

func GetSQSConfig(ctx context.Context, shouldLog bool, prefix, region, account, endpoint string) *Config {
	var queuePrefix *string
	if prefix != "" {
		queuePrefix = &prefix
	}

	return &Config{
		Account:   account,
		Prefix:    queuePrefix,
		shouldLog: shouldLog,
		Region:    region,
		Endpoint:  endpoint,
	}
}

func GetUrl(ctx context.Context, config *Config, name string) (url *string, err error) {
	client, clientErr := getClient(ctx, config)
	if clientErr != nil {
		err = clientErr
		return
	}

	result, urlErr := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &name,
	})
	if urlErr != nil {
		err = urlErr
		return
	}

	url = result.QueueUrl
	return
}

func getClient(ctx context.Context, config *Config) (sqsClient *sqs.SQS, err error) {
	awsConfig := aws.Config{Region: aws.String(config.Region)}
	if config.shouldLog {
		awsConfig.LogLevel = aws.LogLevel(aws.LogDebugWithRequestErrors)
	}

	if endpoint, ok := os.LookupEnv("LOCAL_SQS_ENDPOINT"); ok {
		awsConfig.Endpoint = aws.String(endpoint)
	}

	awsSession, err := session.NewSession(&awsConfig)
	if err != nil {
		return
	}

	sqsClient = sqs.New(awsSession)
	return
}

// CreateQueue We are not using this method, we are creating queue through AWS console
func CreateQueue(ctx context.Context, config *Config, name string, queueType string) (err error) {
	client, err := getClient(ctx, config)
	if err != nil {
		return
	}

	createQueueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(name),
	}
	if QueueType(queueType) == QueueFifo {
		createQueueInput.Attributes = map[string]*string{
			"FifoQueue": aws.String("true"),
		}
	}
	_, err = client.CreateQueue(createQueueInput)
	return
}
