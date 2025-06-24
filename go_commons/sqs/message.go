package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/omniful/go_commons/compression"
	"time"
)

type MessageAttributeValue struct {
	BinaryValue []byte
	DataType    *string
	StringValue *string
}

type Message struct {
	GroupId         string
	Value           []byte
	ReceiptHandle   string
	Attributes      map[string]string
	DeduplicationId string
	DelayDuration   time.Duration
	Headers         map[string]string

	// In case of compression.None, queue compression will work
	Compression compression.Compression
}

func (m *Message) IsDelaySecondsValid() bool {
	if m == nil || m.DelayDuration.Seconds() < 0 || m.DelayDuration.Seconds() > 900 {
		return false
	}

	return true
}

func (m *Message) RetrieveDelaySeconds() *int64 {
	if m == nil {
		return aws.Int64(0)
	}

	return aws.Int64(int64(m.DelayDuration.Seconds()))
}
