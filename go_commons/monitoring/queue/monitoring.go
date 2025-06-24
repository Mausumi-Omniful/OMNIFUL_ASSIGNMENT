package queue

import (
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/util"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// Matrix names
const (
	QueueMessageProduced           = "queue_message_produced_total"
	QueueMessageConsumed           = "queue_message_consumed_total"
	QueueMessageProcessingDuration = "queue_message_processing_duration"
)

const (
	QueueName = "queue_name"
	Delimiter = "_"
)

type MonitoringAttributes map[string]string

type Monitoring interface {
	RecordMessageReceived(queueName string, labels MonitoringAttributes) error
	RecordMessageConsumed(queueName string, labels MonitoringAttributes) error
	RecordProcessingDuration(queueName string, labels MonitoringAttributes, duration time.Duration) error
}

type queueMonitoring struct {
	messageProduced                 *prometheus.CounterVec
	messageProducedLabels           []string
	messageConsumed                 *prometheus.CounterVec
	messageConsumedLabels           []string
	messageProcessingDuration       *prometheus.HistogramVec
	messageProcessingDurationLabels []string
}

type MonitoringConfig struct {
	Prefix                                string
	MessageProducedCustomLabels           []string
	MessageConsumedCustomLabels           []string
	MessageProcessingDurationCustomLabels []string
}

func NewMonitoring(config MonitoringConfig) Monitoring {
	messageProducedLabels := getLabels(config.MessageProducedCustomLabels)
	messageConsumedLabels := getLabels(config.MessageConsumedCustomLabels)
	messageProcessingDurationLabels := getLabels(config.MessageProcessingDurationCustomLabels)

	_queueMonitoring := &queueMonitoring{
		messageProduced: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: getMetricsName(config.Prefix, QueueMessageProduced, Delimiter),
				Help: "Total number of messages produced to the queue.",
			},
			messageProducedLabels,
		),
		messageProducedLabels: messageProducedLabels,

		messageConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: getMetricsName(config.Prefix, QueueMessageConsumed, Delimiter),
				Help: "Total number of messages consumed from the queue.",
			},
			messageConsumedLabels,
		),
		messageConsumedLabels: messageConsumedLabels,

		messageProcessingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: getMetricsName(config.Prefix, QueueMessageProcessingDuration, Delimiter),
				Help: "Time taken to process messages in the queue.",
			},
			messageProcessingDurationLabels,
		),
		messageProcessingDurationLabels: messageProcessingDurationLabels,
	}

	// Register metrics with Prometheus
	prometheus.MustRegister(_queueMonitoring.messageProduced)
	prometheus.MustRegister(_queueMonitoring.messageConsumed)
	prometheus.MustRegister(_queueMonitoring.messageProcessingDuration)

	return _queueMonitoring
}

func (qm *queueMonitoring) RecordMessageReceived(queueName string, attributes MonitoringAttributes) error {
	attributes[QueueName] = queueName
	counter, err := qm.messageProduced.GetMetricWith(
		getPrometheusLabels(
			qm.messageConsumedLabels,
			attributes,
		),
	)
	if err != nil {
		return err
	}

	counter.Inc()

	return nil
}

func (qm *queueMonitoring) RecordMessageConsumed(queueName string, attributes MonitoringAttributes) error {
	attributes[QueueName] = queueName
	counter, err := qm.messageConsumed.GetMetricWith(
		getPrometheusLabels(
			qm.messageConsumedLabels,
			attributes,
		),
	)
	if err != nil {
		return err
	}

	counter.Inc()

	return nil
}

func (qm *queueMonitoring) RecordProcessingDuration(
	queueName string,
	attributes MonitoringAttributes,
	duration time.Duration,
) error {
	attributes[QueueName] = queueName
	histogram, err := qm.messageProcessingDuration.GetMetricWith(
		getPrometheusLabels(
			qm.messageProcessingDurationLabels,
			attributes,
		),
	)
	if err != nil {
		return err
	}

	histogram.Observe(duration.Seconds())

	return nil
}

func getPrometheusLabels(
	labels []string,
	attributes MonitoringAttributes,
) prometheus.Labels {
	result := make(prometheus.Labels)

	for _, label := range labels {
		value, ok := attributes[label]
		if !ok {
			log.Errorf("Queue monitoring, %s must be present", label)
		}

		result[label] = value
	}

	return result
}

func getLabels(extraLabels []string) []string {
	return util.DeduplicateSlice(
		append(defaultLabels(), extraLabels...),
	)
}

func defaultLabels() []string {
	return []string{QueueName}
}

func getMetricsName(prefix, existingName, delimiter string) string {
	if len(prefix) == 0 {
		return existingName
	}

	return prefix + delimiter + existingName
}
