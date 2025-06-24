package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/jpillora/backoff"
	appConfig "github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/pubsub/interceptor"
	"github.com/omniful/go_commons/sqs"
	"time"
)

// ConsumerGroupHandler represents the consumer group.
type ConsumerGroupHandler struct {
	HandlerMap               map[string]pubsub.IPubSubMessageHandler
	Interceptor              interceptor.Interceptor
	DeadLetterQueuePublisher *sqs.Publisher
	Context                  context.Context
	TransactionName          string
	RetryInterval            time.Duration
}

// Handle Generic handler which calls interceptor if registered otherwise calls handler directly.
func (client *ConsumerGroupHandler) Handle(ctx context.Context, message *pubsub.Message) error {
	handler, ok := client.HandlerMap[message.Topic]
	if !ok {
		return errors.New("unable to locate topic handler")
	}

	// No Interceptor present
	if client.Interceptor == nil {
		return handler.Process(ctx, message)
	}

	err := client.Interceptor(ctx, message, func(ctx context.Context, msg *pubsub.Message) error {
		return handler.Process(ctx, msg)
	})

	return err
}

func (client *ConsumerGroupHandler) ParseMessage(message *sarama.ConsumerMessage) (*pubsub.Message, error) {
	headers := make(map[string]string, 0)
	for _, header := range message.Headers {
		headers[string(header.Key)] = string(header.Value)
	}

	msgObject := &pubsub.Message{
		Topic:     message.Topic,
		Value:     message.Value,
		Key:       fmt.Sprintf("%s", message.Key),
		Timestamp: message.Timestamp,
		Headers:   headers,
	}
	return msgObject, nil
}

// Setup is run before consumer start consuming, is normally used to setup things such as database connections
func (client *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (client *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (client *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			if msg == nil {
				continue
			}

			b := &backoff.Backoff{
				Jitter: true,
				Min:    client.RetryInterval,
				Max:    time.Hour,
			}

			// Parse Kafka Message
			parsedMessage, parseErr := client.ParseMessage(msg)
			if parseErr != nil {
				return parseErr
			}

			// Set Config
			ctx, _ := appConfig.SetConfigInContext(client.Context)

			// Set kafka transaction name to identify transaction
			ctx = context.WithValue(ctx, constants.KafkaConsumerTransactionName, client.TransactionName)

			// Set Request ID to trace requests
			ctx = env.SetKafkaRequestID(ctx, parsedMessage.Headers)

			l := log.DefaultLogger()
			l = l.With(
				log.String(constants.HeaderXOmnifulRequestID, env.GetRequestID(ctx)),
			)
			ctx = log.ContextWithLogger(ctx, l)

			// Retrying 4 times in case of error
			for i := 0; i < 4; i++ {
				handlerErr := client.Handle(ctx, parsedMessage)
				if handlerErr != nil {
					log.Errorf("Error occurred while processing message from kafka. Key: %s"+
						"DateCreated = %v, Topic = %s Partition = %d Offset = %d", string(msg.Key), msg.Timestamp, msg.Topic,
						msg.Partition, msg.Offset)

					// Panic Error
					if errors.As(handlerErr, &error2.CustomError{}) {
						client.publishFailedMessageToDeadLetterQueue(ctx, parsedMessage)

						break
					}

					//Retrying with jitter and backoff coefficient
					d := b.Duration()
					time.Sleep(d)

					// On Last Attempt, Send Failed Message to DLQ
					if i == 3 {
						client.publishFailedMessageToDeadLetterQueue(ctx, parsedMessage)
					}

					continue
				}

				// On success, no need to retry
				break
			}

			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (client *ConsumerGroupHandler) publishFailedMessageToDeadLetterQueue(ctx context.Context, msg *pubsub.Message) {
	if client.DeadLetterQueuePublisher == nil {
		return
	}

	publishErr := client.DeadLetterQueuePublisher.Publish(ctx, &sqs.Message{
		Value:      msg.Value,
		Attributes: msg.Headers,
	})
	if publishErr != nil {
		log.Errorf("Error in publishing message to dead letter queue :: %v", publishErr)
	}

	return
}
