package pubsub

import (
	"encoding/json"
	"time"
)

type Message struct {
	Topic     string
	Value     []byte
	Key       string
	Headers   map[string]string
	Timestamp time.Time
}

type MessageType string

// To unmarshal the message into the provided object msg
// dest should always be a pointer
func (m *Message) To(dest interface{}, messageType MessageType) error {
	err := json.Unmarshal(m.Value, dest)
	if err != nil {
		return err
	}

	return nil
}

// NewEventInBytes return the bytes value for the event registry proto
func NewEventInBytes(m interface{}) ([]byte, error) {
	val, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return val, nil
}
