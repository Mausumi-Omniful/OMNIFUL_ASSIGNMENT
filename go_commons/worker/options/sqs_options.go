package options

import "github.com/omniful/go_commons/worker/configs"

type SqsOption func(options *configs.SqsQueueConfig)

func WithSendBatchedMessagesEnabled() SqsOption {
	return func(config *configs.SqsQueueConfig) {
		config.SendBatchedMessages = true
	}
}
