package listener

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	"github.com/omniful/go_commons/worker/configs"
)

type SQSListener struct {
	Handler       sqs.ISqsMessageHandler
	queueConsumer *sqs.Consumer
	Config        configs.SqsQueueConfig
}

func NewSQSListener(handler sqs.ISqsMessageHandler, config configs.SqsQueueConfig) ListenerServer {
	return &SQSListener{
		Handler: handler,
		Config:  config,
	}
}

func (l *SQSListener) Start(ctx context.Context) {
	logTag := "[Workers][SQSListener][Start] "
	log.Info(logTag + "Started")

	log.Debugf(logTag+"Getting Standard Queue %v", l.Config.QueueName)

	queue, err := l.getQueue(ctx)
	if err != nil || queue == nil {
		panic(fmt.Sprintf(
			logTag+"Initialization error. queue: %v, err : %v, publisher: %+v",
			l.Config.QueueName,
			err,
			queue,
		))
	}

	// Creating New Consumer
	consumer, err := sqs.NewConsumer(queue, uint64(l.Config.WorkerCount), l.Config.ConcurrencyPerWorker, l.Handler,
		10, 30, true, l.Config.SendBatchedMessages)
	if err != nil || consumer == nil {
		panic(fmt.Sprintf(
			logTag+"Initialization error. queue: %v, err : %v, publisher: %+v",
			l.Config.QueueName,
			err,
			queue,
		))
	}

	l.queueConsumer = consumer
	consumer.Start(ctx)
}

func (l *SQSListener) Stop() {
	if l.queueConsumer == nil {
		return
	}

	err := l.queueConsumer.Close()
	if err != nil {
		log.Errorf("Error in closing queue consumer :: %v", err.Error())
		return
	}

	return

}

func (l *SQSListener) GetName() string {
	return l.Config.Name
}

func (l *SQSListener) getQueue(ctx context.Context) (queue *sqs.Queue, err error) {
	if l.Config.IsFifo {
		return sqs.NewFifoQueue(ctx, l.Config.QueueName, &sqs.Config{
			Account:  l.Config.Account,
			Endpoint: l.Config.Region,
			Prefix:   &l.Config.Prefix,
			Region:   l.Config.Region,
		})
	}

	return sqs.NewStandardQueue(ctx, l.Config.QueueName, &sqs.Config{
		Account:  l.Config.Account,
		Endpoint: l.Config.Region,
		Prefix:   &l.Config.Prefix,
		Region:   l.Config.Region,
	})
}
