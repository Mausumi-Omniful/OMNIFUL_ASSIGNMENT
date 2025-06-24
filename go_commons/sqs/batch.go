package sqs

import "github.com/aws/aws-sdk-go/service/sqs"

type batchRequest struct {
	entries []*sqs.SendMessageBatchRequestEntry
	size    int
}

// NewBatch creates a new empty batch
func newBatchRequest() *batchRequest {
	return &batchRequest{
		entries: make([]*sqs.SendMessageBatchRequestEntry, 0, maxBatchMessages),
	}
}

// CanAdd checks if a message can be added to the current batch
func (b *batchRequest) CanAdd(messageSize int) bool {
	return len(b.entries) < maxBatchMessages && b.size+messageSize <= maxBatchPayload
}

// Add adds a message entry to the batch
func (b *batchRequest) Add(entry *sqs.SendMessageBatchRequestEntry, messageSize int) {
	b.entries = append(b.entries, entry)
	b.size += messageSize
}
