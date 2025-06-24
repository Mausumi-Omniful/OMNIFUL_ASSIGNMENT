package sqs

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/newrelic"
)

const (
	SentTimestamp  = "SentTimestamp"
	MessageGroupId = "MessageGroupId"
)

type ConsumerWorker struct {
	*Queue
	name              string
	Handler           ISqsMessageHandler
	concurrency       int
	done              chan bool
	maxMessagesCount  int64
	visibilityTimeout int64
	waitTimeSeconds   int64
	sendBatchMessage  bool
}

func NewConsumerWorker(name string, queue *Queue, handler ISqsMessageHandler, concurrency int, done chan bool, maxMessagesCount, visibilityTimeout, waitTimeSeconds int64, sendBatchMessage bool) *ConsumerWorker {
	return &ConsumerWorker{
		Queue:             queue,
		name:              name,
		Handler:           handler,
		concurrency:       concurrency,
		done:              done,
		maxMessagesCount:  maxMessagesCount,
		visibilityTimeout: visibilityTimeout,
		waitTimeSeconds:   waitTimeSeconds,
		sendBatchMessage:  sendBatchMessage,
	}
}

func (c *ConsumerWorker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			cusErr := errors.New(fmt.Sprintf("panic occurred: %+v, stacktrace: %+v", r, string(debug.Stack())))
			newrelic.RecordEvent(newrelic.Error, map[string]interface{}{
				"Err": cusErr.Error()})
		}
		wg.Done()
	}()

	pool := NewPool(c.concurrency, c.Handler, c.name)
	pool.Run(ctx)

	// Ignoring visibility timeout
	receiveMessageInput := &sqs.ReceiveMessageInput{
		AttributeNames:      aws.StringSlice(getAttributeNames()),
		QueueUrl:            c.Url,
		MaxNumberOfMessages: aws.Int64(c.maxMessagesCount),
		WaitTimeSeconds:     aws.Int64(c.waitTimeSeconds),
		MessageAttributeNames: aws.StringSlice([]string{
			constants.AllMessageAttributes,
		}),
	}

	for {
		output, err := c.client.ReceiveMessageWithContext(ctx, receiveMessageInput)
		if ctx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
			break
		}

		if err != nil {
			err = errors.New(fmt.Sprintf("Error SQS message receive error: %+v, queueName:%+v, queueUrl:%+v", err, c.Name, c.Url))
			newrelic.RecordEvent(newrelic.Error, map[string]interface{}{
				"Err": err.Error()})
			time.Sleep(20 * time.Second)
			continue
		}

		messages := make([]Message, 0)

		for _, value := range output.Messages {
			if value.Body == nil {
				err = c.Remove(*value.ReceiptHandle)
				if err != nil {
					continue
				}
				continue
			}

			task := getTask(value)
			if !c.sendBatchMessage {
				pool.AddTask(&[]Message{task})
			} else {
				messages = append(messages, task)
			}
		}

		if c.sendBatchMessage {
			pool.AddTask(&messages)
		}
	}

	pool.Close()
}

func getTask(value *sqs.Message) Message {
	messageAttribute := make(map[string]string, 0)
	headers := make(map[string]string)
	for key, val := range value.MessageAttributes {
		if val.StringValue != nil {
			messageAttribute[key] = *val.StringValue
		}

	}

	for key, val := range value.Attributes {
		if val != nil {
			headers[key] = *val
		}
	}

	return Message{
		Value:         []byte(*value.Body),
		ReceiptHandle: *value.ReceiptHandle,
		Attributes:    messageAttribute,
		Headers:       headers,
	}
}

func getAttributeNames() []string {
	return []string{SentTimestamp, MessageGroupId}
}
