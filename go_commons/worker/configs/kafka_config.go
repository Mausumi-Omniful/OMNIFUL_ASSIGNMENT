package configs

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/config"
)

type KafkaConsumerConfig struct {
	Name              string
	Topic             string
	GroupID           string
	Enabled           bool
	ClientID          string
	KafkaVersion      string
	Brokers           []string
	Region            string
	IAMAuthentication bool
	WorkerGroup       string
}

func GetKafkaConfig(ctx context.Context, consumerName string) KafkaConsumerConfig {
	return KafkaConsumerConfig{
		Name:         config.GetString(ctx, fmt.Sprintf("consumers.%s.name", consumerName)),
		Topic:        config.GetString(ctx, fmt.Sprintf("consumers.%s.topic", consumerName)),
		GroupID:      config.GetString(ctx, fmt.Sprintf("consumers.%s.groupID", consumerName)),
		Enabled:      config.GetBool(ctx, fmt.Sprintf("consumers.%s.enabled", consumerName)),
		ClientID:     config.GetString(ctx, "kafka.clientId"),
		KafkaVersion: config.GetString(ctx, "kafka.version"),
		Brokers:      config.GetStringSlice(ctx, "kafka.brokers"),
		Region:       config.GetString(ctx, "kafka.region"),
		WorkerGroup:  config.GetString(ctx, fmt.Sprintf("consumers.%s.workerGroup", consumerName)),
	}
}
