package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/omniful/go_commons/compression"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"strconv"
)

const (
	// SQS message size limit is 262144 bytes. Reserving 12144 bytes for message attributes
	maxMessageSize  = 250000
	maxBatchPayload = 250000

	// SQS limit of 10 messages per batch
	maxBatchMessages = 10
)

type Publisher struct {
	*Queue
}

func NewPublisher(queue *Queue) *Publisher {
	return &Publisher{
		queue,
	}
}

func (p *Publisher) Publish(ctx context.Context, message *Message) (err error) {
	segment := newrelic.StartSegmentWithContext(ctx, "SQSPublish")
	defer segment.End()

	p.addCompressionToMessage(message)
	return p.send(ctx, message)
}

func (p *Publisher) BatchPublish(ctx context.Context, messages []*Message) (err error) {
	segment := newrelic.StartSegmentWithContext(ctx, "SQSBatchPublish")
	defer segment.End()

	if len(messages) == 0 {
		return
	}

	for _, message := range messages {
		p.addCompressionToMessage(message)
	}

	return p.sendBatch(ctx, messages)
}

func (p *Publisher) addCompressionToMessage(msg *Message) {
	if len(msg.Value) > maxMessageSize {
		msg.Compression = compression.GZIP
	}
}

func (p *Publisher) send(ctx context.Context, message *Message) (err error) {
	// Compress message and get attributes
	msgBody, msgAttrs, err := p.compressMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("message preparation failed: %w", err)
	}

	sendMessageRequest := &sqs.SendMessageInput{
		MessageBody:       aws.String(string(msgBody)),
		QueueUrl:          p.Url,
		MessageAttributes: msgAttrs,
	}

	if message.IsDelaySecondsValid() {
		sendMessageRequest.DelaySeconds = message.RetrieveDelaySeconds()
	}

	if p.Type == QueueFifo {
		sendMessageRequest.MessageGroupId = aws.String(message.GroupId)
	}

	if len(message.DeduplicationId) > 0 {
		sendMessageRequest.MessageDeduplicationId = aws.String(message.DeduplicationId)
	}

	_, err = p.client.SendMessage(sendMessageRequest)
	return
}

func (p *Publisher) compressMessage(ctx context.Context, msg *Message) ([]byte, map[string]*sqs.MessageAttributeValue, error) {
	attrs := ParseMessageAttributes(ctx, msg.Attributes)

	compressor := p.getCompressor(msg)

	// If no compression specified, return original message
	if compressor.Compression() == compression.None {
		return msg.Value, attrs, nil
	}

	compressedMsg, err := compressor.Compress(msg.Value)

	// Handle successful compression
	if err == nil {
		return compressedMsg, addCompressionAttribute(attrs, compressor.Compression()), nil
	}

	log.Errorf("Message compression failed, using compression ::%v : %v", compressor.Compression(), err)
	newrelic.NoticeError(ctx, err)
	return msg.Value, attrs, nil
}

func (p *Publisher) getCompressor(msg *Message) compression.Compressor {
	if msg.Compression != compression.None {
		return compression.GetCompressionParser(msg.Compression)
	}

	return p.compressor
}

func addCompressionAttribute(attrs map[string]*sqs.MessageAttributeValue, compression compression.Compression) map[string]*sqs.MessageAttributeValue {
	attrs[constants.Compression] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(strconv.Itoa(int(compression))),
	}
	return attrs
}

func ParseMessageAttributes(ctx context.Context, attributes map[string]string) map[string]*sqs.MessageAttributeValue {
	if attributes == nil {
		attributes = make(map[string]string)
	}

	// Add or update the XRequestID attribute if not present
	if _, ok := attributes[constants.HeaderXOmnifulRequestID]; !ok {
		attributes[constants.HeaderXOmnifulRequestID] = env.GetRequestID(ctx)
	}

	sqsAttributeMap := make(map[string]*sqs.MessageAttributeValue)
	for attributeKey, attributeValue := range attributes {
		sqsAttributeMap[attributeKey] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attributeValue),
		}
	}

	return sqsAttributeMap
}

func (p *Publisher) sendBatch(ctx context.Context, messages []*Message) error {
	batches, batchErr := p.createBatches(ctx, messages)
	if batchErr != nil {
		return fmt.Errorf("failed to create batches: %w", batchErr)
	}

	for _, batch := range batches {
		if sendErr := p.sendMessageBatch(batch); sendErr != nil {
			return fmt.Errorf("failed to send message batch: %w", sendErr)
		}
	}
	return nil
}

func (p *Publisher) createBatches(ctx context.Context, messages []*Message) ([]*batchRequest, error) {
	var batches []*batchRequest
	currentBatch := newBatchRequest()

	for i, msg := range messages {
		msgBody, msgAttrs, err := p.compressMessage(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("message compression failed for message %d: %w", i, err)
		}

		messageSize := len(msgBody)
		if messageSize > maxMessageSize {
			return nil, fmt.Errorf("message %d exceeds maximum size of %d bytes after compression: got %d bytes",
				i, maxMessageSize, messageSize)
		}

		entry := &sqs.SendMessageBatchRequestEntry{
			Id:                aws.String(strconv.Itoa(i)),
			MessageBody:       aws.String(string(msgBody)),
			MessageAttributes: msgAttrs,
		}

		if p.Type == QueueFifo {
			entry.MessageGroupId = aws.String(msg.GroupId)
		}

		if !currentBatch.CanAdd(messageSize) {
			batches = append(batches, currentBatch)
			currentBatch = newBatchRequest()
		}

		currentBatch.Add(entry, messageSize)
	}

	batches = append(batches, currentBatch)
	return batches, nil
}

func (p *Publisher) sendMessageBatch(batch *batchRequest) error {
	_, err := p.client.SendMessageBatch(&sqs.SendMessageBatchInput{
		QueueUrl: p.Url,
		Entries:  batch.entries,
	})
	return err
}
