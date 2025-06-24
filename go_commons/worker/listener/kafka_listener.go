package listener

import (
	"context"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/worker/configs"
)

type KafkaListener struct {
	Handler  pubsub.IPubSubMessageHandler
	consumer *kafka.ConsumerClient
	Config   configs.KafkaConsumerConfig
}

func (l *KafkaListener) Start(ctx context.Context) {
	logTag := "[Workers][KafkaListener][Start] "
	log.Info(logTag + "Started")

	if !l.Config.Enabled {
		log.Infof(logTag+"kafka consumer of Hub Listener not working for topic %s", l.Config.Topic)
		return
	}

	l.consumer = kafka.NewConsumer(
		kafka.WithBrokers(l.Config.Brokers),
		kafka.WithClientID(l.Config.ClientID),
		kafka.WithKafkaVersion(l.Config.KafkaVersion),
		kafka.WithConsumerGroup(l.Config.GroupID),
		kafka.WithRegion(l.Config.Region),
		kafka.WithIAMAuthentication(l.Config.IAMAuthentication),
	)

	l.consumer.RegisterHandler(l.Config.Topic, l.Handler)

	log.Printf(logTag+"starting kafka consumer for topic %s, consumer_group: %s", l.Config.Topic, l.Config.GroupID)
	l.consumer.Subscribe(ctx)
}

func (l *KafkaListener) Stop() {
	if l.consumer != nil {
		l.consumer.Close()
	}
}

func (l *KafkaListener) GetName() string {
	return l.Config.Name
}

func NewKafkaListener(handler pubsub.IPubSubMessageHandler, config configs.KafkaConsumerConfig) ListenerServer {
	return &KafkaListener{
		Handler: handler,
		Config:  config,
	}
}
