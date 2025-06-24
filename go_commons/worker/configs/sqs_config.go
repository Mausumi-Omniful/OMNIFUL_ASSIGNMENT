package configs

import (
	"context"
	"fmt"
	"strings"

	"github.com/omniful/go_commons/config"
)

type SqsQueueConfig struct {
	QueueName            string
	WorkerCount          int64
	WorkerEnabled        bool
	Account              string
	Prefix               string
	Region               string
	Endpoint             string
	ShouldLog            bool
	ConcurrencyPerWorker uint64
	Name                 string
	IsFifo               bool
	WorkerGroup          string
	SendBatchedMessages  bool
}

func GetSqsConfig(ctx context.Context, consumerName string) SqsQueueConfig {
	queueName := config.GetString(ctx, fmt.Sprintf("workers.%s.name", consumerName))
	return SqsQueueConfig{
		QueueName:            queueName,
		WorkerCount:          config.GetInt64(ctx, fmt.Sprintf("workers.%s.workerCount", consumerName)),
		WorkerEnabled:        config.GetBool(ctx, fmt.Sprintf("workers.%s.workerEnabled", consumerName)),
		Account:              config.GetString(ctx, "aws.account"),
		Region:               config.GetString(ctx, "aws.region"),
		ShouldLog:            config.GetBool(ctx, "aws.shouldLog"),
		Prefix:               config.GetString(ctx, "aws.sqs.prefix"),
		ConcurrencyPerWorker: config.GetUint64(ctx, fmt.Sprintf("workers.%s.concurrencyPerWorker", consumerName)),
		Name:                 config.GetString(ctx, fmt.Sprintf("workers.%s.workerName", consumerName)),
		IsFifo:               strings.HasSuffix(queueName, ".fifo"),
		WorkerGroup:          config.GetString(ctx, fmt.Sprintf("workers.%s.workerGroup", consumerName)),
	}
}
