package sqs

import (
	"context"
	"github.com/omniful/go_commons/compression"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"strconv"
)

type HandlerWrapper struct {
	handler ISqsMessageHandler
	queue   *Queue
}

func NewHandlerWrapper(handler ISqsMessageHandler, queue *Queue) *HandlerWrapper {
	return &HandlerWrapper{
		handler: handler,
		queue:   queue,
	}
}

func (h *HandlerWrapper) Process(ctx context.Context, messages *[]Message) error {
	txn := newrelic.StartTransaction(h.queue.Name, nil, nil)
	ctx = newrelic.NewContext(context.TODO(), txn)
	defer txn.End()

	l := log.DefaultLogger()

	if len(*messages) == 1 {
		ctx = env.SetSqsMessageRequestID(ctx, (*messages)[0].Attributes)
		requestID := env.GetRequestID(ctx)

		newrelic.SetRequestID(ctx, requestID)
		l = l.With(
			log.String(constants.HeaderXOmnifulRequestID, requestID),
		)
	}

	ctx, _ = config.SetConfigInContext(ctx)
	ctx = log.ContextWithLogger(ctx, l)

	// Decompress the messages
	for i := 0; i < len(*messages); i++ {
		attributes := (*messages)[i].Attributes

		c, ok := attributes[constants.Compression]
		if !ok {
			continue
		}

		comp, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			txn.NoticeError(err)
			return err
		}

		parser := compression.GetCompressionParser(compression.Compression(int8(comp)))

		decompressMsg, err := parser.Decompress((*messages)[i].Value)
		if err != nil {
			txn.NoticeError(err)
			return err
		}

		(*messages)[i].Value = decompressMsg
	}

	err := h.handler.Process(ctx, messages)
	if err != nil {
		txn.NoticeError(err)
		return err
	}

	for _, message := range *messages {
		err = h.queue.Remove(message.ReceiptHandle)
		if err != nil {
			return err
		}
	}
	return nil
}
