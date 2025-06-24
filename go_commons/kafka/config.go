package kafka

import (
	"github.com/IBM/sarama"
	"github.com/omniful/go_commons/compression"
	"github.com/omniful/go_commons/log"
	"time"
)

type config struct {
	saramaConfig     *sarama.Config
	brokers          []string
	consumerGroupID  string
	transactionName  string
	retryInterval    time.Duration
	clientID         string
	version          string
	region           string
	useSASLProtocol  bool
	userName         string
	password         string
	saslMechanism    sarama.SASLMechanism
	deadLetterConfig *DeadLetterQueueConfig
	compression      compression.Compression
	asyncProducer    bool
}

func (c *config) getTransactionName() string {
	if c.transactionName != "" {
		return c.transactionName
	}

	return c.consumerGroupID
}

func (c *config) getCompression() sarama.CompressionCodec {
	switch c.compression {
	case compression.GZIP:
		return sarama.CompressionGZIP
	default:
		return sarama.CompressionSnappy
	}
}

type DeadLetterQueueConfig struct {
	Queue     string
	Account   string
	Endpoint  string
	Prefix    string
	ShouldLog bool
	Region    string
}

type option func(*config)

// WithBrokers function sets the number of kafka brokers
func WithBrokers(brokers []string) func(*config) {
	return func(conf *config) {
		conf.brokers = brokers
	}
}

// WithConsumerGroup sets group name for kafka consumer group
func WithConsumerGroup(group string) func(*config) {
	return func(conf *config) {
		conf.consumerGroupID = group
	}
}

// WithTransactionName sets transaction name for monitoring tools
func WithTransactionName(name string) func(*config) {
	return func(conf *config) {
		conf.transactionName = name
	}
}

// WithRetryInterval sets retry time interval for retrying message
func WithRetryInterval(duration time.Duration) func(*config) {
	return func(conf *config) {
		conf.retryInterval = duration
	}
}

// WithClientID takes the client id to use for creating sarama config
func WithClientID(clientID string) func(*config) {
	return func(conf *config) {
		conf.clientID = clientID
	}
}

// WithAsyncProducer takes the that asyncProducer is needed or not
func WithAsyncProducer(asyncProducer bool) func(*config) {
	return func(conf *config) {
		conf.asyncProducer = asyncProducer
	}
}

// WithKafkaVersion takes the kafka version to use for creating sarama config
func WithKafkaVersion(version string) func(*config) {
	return func(conf *config) {
		conf.version = version
	}
}

// WithDeadLetterConfig function sets dead-letter queue
func WithDeadLetterConfig(deadLetterConfig *DeadLetterQueueConfig) func(*config) {
	return func(conf *config) {
		conf.deadLetterConfig = deadLetterConfig
	}
}

// WithRegion function sets the region for IAM authentication
func WithRegion(region string) func(*config) {
	return func(conf *config) {
		conf.region = region
	}
}

// WithIAMAuthentication function sets the IAM authentication for kafka cluster
func WithIAMAuthentication(enable bool) func(*config) {
	return func(conf *config) {
		conf.useSASLProtocol = enable
		conf.saslMechanism = sarama.SASLTypeOAuth
	}
}

// WithSASLPlainAuthentication function sets the SASL Plain Credentials for kafka cluster
func WithSASLPlainAuthentication(user, password string) func(*config) {
	return func(conf *config) {
		conf.useSASLProtocol = true
		conf.userName = user
		conf.password = password
		conf.saslMechanism = sarama.SASLTypePlaintext
	}
}

// WithDataCompression function sets the compression for compressing data
func WithDataCompression(compression compression.Compression) func(*config) {
	return func(conf *config) {
		conf.compression = compression
	}
}

func parseKafkaVersion(version string) sarama.KafkaVersion {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Not a valid kafka version %s", err)
	}

	return kafkaVersion
}
