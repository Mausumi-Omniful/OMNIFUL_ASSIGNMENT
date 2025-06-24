package kafka

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"

	"github.com/IBM/sarama"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/pubsub/interceptor"
)

// ConsumerClient type represents sarama consumer group interface
type ConsumerClient struct {
	config                   *config
	consumer                 sarama.ConsumerGroup
	HandlerMap               map[string]pubsub.IPubSubMessageHandler
	Interceptor              interceptor.Interceptor
	DeadLetterQueuePublisher *sqs.Publisher
	transactionName          string
	mutex                    sync.Mutex
}

func NewConsumer(opts ...option) *ConsumerClient {
	config := &config{
		retryInterval: time.Second, // Default Retry Interval
	}

	// Override config
	for _, opt := range opts {
		opt(config)
	}

	log.Debugf("kafka consumer config debug, %v", config)

	consumer := newConsumerGroup(config)
	deadLetterQueuePublisher := newDeadLetterQueuePublisher(config)

	consumerClient := &ConsumerClient{
		consumer:                 consumer,
		config:                   config,
		HandlerMap:               map[string]pubsub.IPubSubMessageHandler{},
		Interceptor:              interceptor.NewRelicInterceptor(),
		DeadLetterQueuePublisher: deadLetterQueuePublisher,
		transactionName:          config.getTransactionName(),
	}
	return consumerClient
}

func newConsumerGroup(config *config) sarama.ConsumerGroup {
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
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Fetch.Min = 50 * 1024
	saramaConfig.Consumer.MaxWaitTime = 400 * time.Millisecond
	saramaConfig.Metadata.RefreshFrequency = 60 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(config.brokers, config.consumerGroupID, saramaConfig)

	if err != nil {
		log.Panicf("Failed to start Sarama ConsumerGroup: %v", err)
	}

	return consumerGroup
}

func newDeadLetterQueuePublisher(config *config) *sqs.Publisher {
	if config.deadLetterConfig == nil {
		return nil
	}

	queue, sqsErr := sqs.NewStandardQueue(context.TODO(), config.deadLetterConfig.Queue, &sqs.Config{
		Account:  config.deadLetterConfig.Account,
		Endpoint: config.deadLetterConfig.Endpoint,
		Region:   config.deadLetterConfig.Region,
		Prefix:   aws.String(config.deadLetterConfig.Prefix),
	})
	if sqsErr != nil || queue == nil {
		log.Panicf("Failed to start Dead Letter Queue: %v", config.deadLetterConfig.Queue)
		return nil
	}

	return sqs.NewPublisher(queue)
}

// Close function stops kafka consumer to listen from any new message
func (client *ConsumerClient) Close() {
	if client.consumer != nil {
		log.Infof("Closing kafka Consumer")
		err := client.consumer.Close()
		if err != nil {
			log.Error("Unable to close sarama consumer", err)
		}
	}
}

func (client *ConsumerClient) SetInterceptor(interceptor interceptor.Interceptor) *ConsumerClient {
	if interceptor != nil {
		client.Interceptor = interceptor
	}
	return client
}

// RegisterHandler registers a handler for a topic in a concurrent-safe way
func (client *ConsumerClient) RegisterHandler(topic string, handler pubsub.IPubSubMessageHandler) *ConsumerClient {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.HandlerMap[topic] = handler
	return client
}

// UnRegisterHandler removes a handler for a topic in a concurrent-safe way
func (client *ConsumerClient) UnRegisterHandler(topic string) *ConsumerClient {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	delete(client.HandlerMap, topic)
	return client
}

func (client *ConsumerClient) Subscribe(ctx context.Context) {
	go func() {
		for err := range client.consumer.Errors() {
			if consumerErr, ok := err.(*sarama.ConsumerError); ok {
				log.Error("KafkaErrorStream: ", consumerErr)
			} else {
				log.Error("KafkaErrorStreamUnknownError: ", err)
			}
		}
	}()

	for {
		if len(client.HandlerMap) == 0 {
			log.Panicf("HandlerMap is empty. please register handlers for topics")
			return
		}

		consumerHandler := &ConsumerGroupHandler{
			HandlerMap:               client.HandlerMap,
			Interceptor:              client.Interceptor,
			DeadLetterQueuePublisher: client.DeadLetterQueuePublisher,
			Context:                  ctx,
			TransactionName:          client.transactionName,
			RetryInterval:            client.config.retryInterval,
		}

		topics := make([]string, 0, 0)
		for topic := range client.HandlerMap {
			topics = append(topics, topic)
		}

		err := client.consumer.Consume(ctx, topics, consumerHandler)
		if err != nil && !errors.Is(err, sarama.ErrClosedConsumerGroup) {
			log.Errorf("KafkaConsumeError: %v", err)
		}

		// Consume will exit in case of rebalance
		time.Sleep(1 * time.Second)
	}
}
