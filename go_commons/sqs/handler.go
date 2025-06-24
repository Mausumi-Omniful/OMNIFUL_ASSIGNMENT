package sqs

import (
	"context"
)

type ISqsMessageHandler interface {
	Process(ctx context.Context, message *[]Message) error
}
