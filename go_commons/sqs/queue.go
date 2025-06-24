package sqs

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/omniful/go_commons/compression"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Account   string
	Endpoint  string
	Prefix    *string
	shouldLog bool
	Region    string

	// Defaults to compression.None
	Compression compression.Compression
}

type Queue struct {
	Name   string
	Url    *string
	Type   QueueType
	client *sqs.SQS
	*Config
	compressor compression.Compressor
}

func NewStandardQueue(ctx context.Context, name string, config *Config) (queue *Queue, err error) {
	if strings.HasSuffix(name, ".fifo") {
		err = errors.New("Invalid standard queue name:" + name)
		return
	}

	queue, err = newQueue(ctx, name, QueueStandard, config)
	return
}

func NewFifoQueue(ctx context.Context, name string, config *Config) (queue *Queue, err error) {
	if strings.HasSuffix(name, ".fifo") {
		queue, err = newQueue(ctx, name, QueueFifo, config)
	} else {
		err = errors.New("Invalid fifo queue name:" + name)
	}

	return
}

func newQueue(ctx context.Context, name string, queueType QueueType, config *Config) (queue *Queue, err error) {
	if shouldLog, ok := os.LookupEnv("AWS_DEBUG_LOG"); ok {
		config.shouldLog, _ = strconv.ParseBool(shouldLog)
	}

	client, err := getClient(ctx, config)
	if err != nil {
		return nil, err
	}

	result, urlErr := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if urlErr != nil {
		err = urlErr
		return nil, err
	}

	compressor := compression.GetCompressionParser(config.Compression)

	queue = &Queue{
		Url:        result.QueueUrl,
		client:     client,
		Name:       name,
		Type:       queueType,
		compressor: compressor,
		Config:     config,
	}

	return queue, nil
}
