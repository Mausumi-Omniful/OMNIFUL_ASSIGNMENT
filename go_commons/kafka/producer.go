package kafka

import (
	"context"
	"errors"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"time"

	"github.com/IBM/sarama"
	"github.com/omniful/go_commons/pubsub"
)

const (
	DialTimeout  = 10 * time.Second
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
)

type ProducerClient struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	config        *config
}

func NewProducer(opts ...option) *ProducerClient {
	c := &config{}
	// Override config
	for _, opt := range opts {
		opt(c)
	}

	producerClient := &ProducerClient{
		syncProducer: newSyncProducer(c),
		config:       c,
	}

	if c.asyncProducer {
		producerClient.asyncProducer = newAsyncProducer(c)
	}

	log.Infof("kafka producer config debug, %v", c)
	return producerClient
}

func newSyncProducer(config *config) sarama.SyncProducer {
	saramaConfig := sarama.NewConfig()
	if config.useSASLProtocol && (config.saslMechanism == sarama.SASLTypeOAuth || config.saslMechanism == sarama.SASLTypePlaintext) {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.SASL.Mechanism = config.saslMechanism

		if config.saslMechanism == sarama.SASLTypeOAuth {
			saramaConfig.Net.SASL.TokenProvider = &MSKAccessTokenProvider{
				region: config.region,
			}
		}

		if config.saslMechanism == sarama.SASLTypePlaintext {
			saramaConfig.Net.SASL.User = config.userName
			saramaConfig.Net.SASL.Password = config.password
		}
	}

	saramaConfig.ClientID = config.clientID
	saramaConfig.Version = parseKafkaVersion(config.version)
	saramaConfig.Net.ReadTimeout = ReadTimeout
	saramaConfig.Net.WriteTimeout = WriteTimeout
	saramaConfig.Net.DialTimeout = DialTimeout
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll             // Wait for all in-sync replicas to ack the message
	saramaConfig.Producer.Compression = config.getCompression()        // Compress Message
	saramaConfig.Producer.MaxMessageBytes = int(sarama.MaxRequestSize) // Refer: https://github.com/IBM/sarama/issues/2142
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Retry.Backoff = time.Second
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Metadata.RefreshFrequency = 60 * time.Second

	producer, err := sarama.NewSyncProducer(config.brokers, saramaConfig)
	if err != nil {
		log.Panicf("Failed to start Sarama SyncProducer: %v", err)
	}

	return producer
}

func newAsyncProducer(config *config) sarama.AsyncProducer {
	saramaConfig := sarama.NewConfig()
	if config.useSASLProtocol && (config.saslMechanism == sarama.SASLTypeOAuth || config.saslMechanism == sarama.SASLTypePlaintext) {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.SASL.Mechanism = config.saslMechanism

		if config.saslMechanism == sarama.SASLTypeOAuth {
			saramaConfig.Net.SASL.TokenProvider = &MSKAccessTokenProvider{
				region: config.region,
			}
		}

		if config.saslMechanism == sarama.SASLTypePlaintext {
			saramaConfig.Net.SASL.User = config.userName
			saramaConfig.Net.SASL.Password = config.password
		}
	}

	saramaConfig.ClientID = config.clientID
	saramaConfig.Version = parseKafkaVersion(config.version)
	saramaConfig.Net.ReadTimeout = ReadTimeout
	saramaConfig.Net.WriteTimeout = WriteTimeout
	saramaConfig.Net.DialTimeout = DialTimeout
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal           // Only wait for the leader to ack
	saramaConfig.Producer.Compression = config.getCompression()        // Compress Message
	saramaConfig.Producer.MaxMessageBytes = int(sarama.MaxRequestSize) // Refer: https://github.com/IBM/sarama/issues/214
	saramaConfig.Producer.Flush.Frequency = 500 * time.Millisecond     // Flush batches every 500ms
	saramaConfig.Producer.Return.Errors = false                        // if setting this to true, ensure to drain errors.
	saramaConfig.Metadata.RefreshFrequency = 60 * time.Second

	producer, err := sarama.NewAsyncProducer(config.brokers, saramaConfig)
	if err != nil {
		log.Panicf("Failed to start sarama AsyncProducer: %s", err)
	}

	return producer
}

func (client *ProducerClient) Publish(ctx context.Context, msg *pubsub.Message) error {
	//Validator segment
	segment := newrelic.StartSegmentWithContext(ctx, "KafkaPublish")
	defer segment.End()

	// Adding request ID
	if len(msg.Headers) == 0 {
		msg.Headers = make(map[string]string)
	}

	msg.Headers[constants.HeaderXOmnifulRequestID] = env.GetRequestID(ctx)

	if client.syncProducer != nil {
		_, _, err := client.syncProducer.SendMessage(client.buildMessage(msg))
		return err
	} else {
		log.Panicf("sync producer not initialized")
	}

	return nil
}

func (client *ProducerClient) PublishBatch(ctx context.Context, msgs []*pubsub.Message) error {
	segment := newrelic.StartSegmentWithContext(ctx, "KafkaBatchPublish")
	defer segment.End()

	messages := make([]*sarama.ProducerMessage, 0)
	for _, msg := range msgs {
		// Adding request ID
		if len(msg.Headers) == 0 {
			msg.Headers = make(map[string]string)
		}

		msg.Headers[constants.HeaderXOmnifulRequestID] = env.GetRequestID(ctx)
		messages = append(messages, client.buildMessage(msg))
	}

	if client.syncProducer != nil {
		err := client.syncProducer.SendMessages(messages)
		return err
	} else {
		log.Panicf("sync producer not initialized")
	}

	return nil
}

func (client *ProducerClient) PublishAsync(ctx context.Context, msg *pubsub.Message) error {
	if !client.config.asyncProducer {
		return errors.New("async producer not initialised")
	}

	// Adding request ID
	if len(msg.Headers) == 0 {
		msg.Headers = make(map[string]string)
	}

	// Adding request ID
	msg.Headers[constants.HeaderXOmnifulRequestID] = env.GetRequestID(ctx)

	if client.asyncProducer != nil {
		client.asyncProducer.Input() <- client.buildMessage(msg)
	} else {
		log.Panicf("async producer not initialized")
	}
	return nil
}

// Close : For releasing producer resources
func (client *ProducerClient) Close() {
	if client.syncProducer != nil {
		err := client.syncProducer.Close()
		if err != nil {
			log.Panicf("Unable to close sync producer", err)
		}
	}

	if client.asyncProducer != nil {
		err := client.asyncProducer.Close()
		if err != nil {
			log.Panicf("unable to close async producer : %s", err.Error())
		}
	}
}

func (client *ProducerClient) Receive(process func(message *pubsub.Message) error) error {
	return nil
}

func (client *ProducerClient) Remove(receiptHandle string) error {
	return nil
}

func (client *ProducerClient) buildMessage(msg *pubsub.Message) *sarama.ProducerMessage {
	kafkaHeaders := make([]sarama.RecordHeader, 0)
	for key, val := range msg.Headers {
		kafkaHeaders = append(kafkaHeaders, sarama.RecordHeader{
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(val),
		})
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic:   msg.Topic,
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: kafkaHeaders,
	}

	if msg.Key != "" {
		kafkaMessage.Key = sarama.StringEncoder(msg.Key)
	}

	return kafkaMessage
}
