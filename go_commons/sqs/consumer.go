package sqs

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/newrelic"
	"github.com/omniful/go_commons/shutdown"
	"sync"
)

type Consumer struct {
	*Queue
	maxMessagesCount     int64
	visibilityTimeout    int64
	waitTimeSeconds      int64
	isAsync              bool
	Handler              ISqsMessageHandler
	done                 chan bool
	numberOfWorker       uint64
	concurrencyPerWorker uint64
	wg                   sync.WaitGroup
	sendBatchMessage     bool
	cancelFunc           context.CancelFunc
}

func NewConsumer(
	queue *Queue,
	numberOfWorker uint64,
	concurrencyPerWorker uint64,
	handler ISqsMessageHandler,
	maxMessagesCount int64,
	visibilityTimeout int64,
	isAsync bool,
	sendBatchMessage bool,
) (*Consumer, error) {
	if maxMessagesCount > 10 {
		return nil, errors.New("maxMessagesCount can not be greater than 10")
	}

	// Set concurrencyPerWorker equal to one and batch to false in case of fifo queue
	if queue.Type == QueueFifo {
		sendBatchMessage = false
		maxMessagesCount = 1
	}

	// Set numberOfWorker and concurrencyPerWorker equal to one in case of sync
	if !isAsync {
		numberOfWorker = 1
		concurrencyPerWorker = 1
	}

	// Worker count should not be more than 2
	if numberOfWorker > 2 {
		numberOfWorker = 2
	}

	// waitingTime is 20
	return &Consumer{
		Queue:                queue,
		numberOfWorker:       numberOfWorker,
		concurrencyPerWorker: concurrencyPerWorker,
		maxMessagesCount:     maxMessagesCount,
		visibilityTimeout:    visibilityTimeout,
		waitTimeSeconds:      20,
		isAsync:              isAsync,
		Handler:              NewHandlerWrapper(handler, queue),
		done:                 make(chan bool, numberOfWorker),
		sendBatchMessage:     sendBatchMessage,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	// Register shutdown callback
	shutdown.RegisterShutdownCallback(c.Name, c)
	ctx, c.cancelFunc = context.WithCancel(ctx)

	for i := 1; i <= int(c.numberOfWorker); i++ {
		id := uuid.New().String()
		consumerWorker := NewConsumerWorker(c.Name+id, c.Queue, c.Handler, int(c.concurrencyPerWorker), c.done, c.maxMessagesCount, c.visibilityTimeout, c.waitTimeSeconds, c.sendBatchMessage)
		c.wg.Add(1)
		go consumerWorker.Start(ctx, &c.wg)
	}
}

func (queue *Queue) Remove(receiptHandle string) (err error) {
	deleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      queue.Url,
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, err = queue.client.DeleteMessage(deleteMessageInput)
	if err != nil {
		newrelic.RecordEvent(newrelic.Error, map[string]interface{}{
			"Err": fmt.Sprintf("Error Unable to delete message from queue deleteMessageInput: %+v. error: %+v", deleteMessageInput, err.Error())})
	}
	return
}

func (c *Consumer) Close() error {
	c.cancelFunc()
	c.wg.Wait()
	return nil
}
