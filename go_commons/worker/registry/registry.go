package registry

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/worker/options"

	commonhttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/pubsub"
	"github.com/omniful/go_commons/sqs"
	"github.com/omniful/go_commons/util"
	"github.com/omniful/go_commons/worker/configs"
	"github.com/omniful/go_commons/worker/listener"
)

// KafkaHandlerInitializer is a type definition for Kafka listener initializers
type KafkaHandlerInitializer func(ctx context.Context, config configs.KafkaConsumerConfig) pubsub.IPubSubMessageHandler

// SQSHandlerInitializer is a type definition for Sqs listener initializers
type SQSHandlerInitializer func(ctx context.Context, config configs.SqsQueueConfig) sqs.ISqsMessageHandler

// Registry manages a collection of listeners grouped by their configuration settings.
type Registry struct {
	listeners            map[string]listener.Configs
	defaultListeners     []listener.Config
	listenerConfigByName map[string]listener.Config
}

// NewRegistry creates and returns a new Registry instance with an initialized listeners map.
func NewRegistry() *Registry {
	return &Registry{
		listeners:            make(map[string]listener.Configs),
		defaultListeners:     make([]listener.Config, 0),
		listenerConfigByName: make(map[string]listener.Config),
	}
}

// RegisterKafkaListenerConfig creates and registers a Kafka listener configuration
// @param ctx context.Context
// @param consumerName string - name from the config.yaml file (consumers.consumerName)
// @param handlerInitializer KafkaHandlerInitializer
func (r *Registry) RegisterKafkaListenerConfig(
	ctx context.Context,
	consumerName string,
	handlerInitializer KafkaHandlerInitializer,
) {
	kafkaConsumerConfig := configs.GetKafkaConfig(ctx, consumerName)
	listenerConfig := listener.Config{
		Name: kafkaConsumerConfig.Name,
		Initializer: func(ctx context.Context) listener.ListenerServer {
			return listener.NewKafkaListener(handlerInitializer(ctx, kafkaConsumerConfig), kafkaConsumerConfig)
		},
	}

	r.RegisterListener(kafkaConsumerConfig.WorkerGroup, listenerConfig)
}

// RegisterSQSListenerConfig creates and registers an SQS listener configuration
// @param ctx context.Context
// @param consumerName string - name from the config.yaml file (workers.consumerName)
// @param handlerInitializer SQSHandlerInitializer
func (r *Registry) RegisterSQSListenerConfig(
	ctx context.Context,
	consumerName string,
	handlerInitializer SQSHandlerInitializer,
	options ...options.SqsOption,
) {
	sqsConsumerConfig := configs.GetSqsConfig(ctx, consumerName)

	for _, opt := range options {
		opt(&sqsConsumerConfig)
	}

	listenerConfig := listener.Config{
		Name: sqsConsumerConfig.Name,
		Initializer: func(ctx context.Context) listener.ListenerServer {
			return listener.NewSQSListener(handlerInitializer(ctx, sqsConsumerConfig), sqsConsumerConfig)
		},
	}

	r.RegisterListener(sqsConsumerConfig.WorkerGroup, listenerConfig)
}

// RegisterHTTPListenerConfig creates and registers an HTTP server as a default listener
func (r *Registry) RegisterHTTPListenerConfig(
	httpServer *commonhttp.Server,
	serviceName string,
) {
	listenerConfig := listener.Config{
		Name: serviceName,
		Initializer: func(ctx context.Context) listener.ListenerServer {
			return listener.NewHttpListener(httpServer, serviceName)
		},
	}

	r.AddDefaultListener(listenerConfig) // Added default listener instead of group
}

// RegisterListener adds a new listener configuration to the specified group.
// If the group doesn't exist, it will be created.
func (r *Registry) RegisterListener(
	group string,
	listenerConfig listener.Config,
) {
	err := listenerConfig.Validate()
	if err != nil {
		panic(err)
	}

	if _, exists := r.listenerConfigByName[listenerConfig.Name]; exists {
		panic(fmt.Sprintf("Listener with name %s already registered", listenerConfig.Name))
	}

	if _, exists := r.listeners[group]; !exists {
		r.listeners[group] = make(listener.Configs, 0)
	}

	r.listenerConfigByName[listenerConfig.Name] = listenerConfig
	r.listeners[group] = append(r.listeners[group], listenerConfig)
}

// AddDefaultListener adds a new listener configuration to the specified group.
// If the group doesn't exist, it will be created.
func (r *Registry) AddDefaultListener(
	listenerConfig listener.Config,
) {
	err := listenerConfig.Validate()
	if err != nil {
		panic(err)
	}

	r.defaultListeners = append(r.defaultListeners, listenerConfig)
}

// GetListenersFromConfig returns a list of listener servers based on the provided configuration.
// It filters listeners based on included/excluded groups and specific listener names.
// If a specific listener name is provided in the config, only that listener is returned.
func (r *Registry) GetListenersFromConfig(
	ctx context.Context,
	config configs.ServerConfig,
) listener.ListenerServers {
	listeners := r.getDefaultListeners(ctx)
	if len(config.GetListenerNames()) > 0 {
		for _, listenerName := range config.GetListenerNames() {
			listeners = append(listeners, r.getListenerFromName(ctx, listenerName))
		}

		return listeners
	}

	groups := r.getAllGroups()
	if len(config.GetIncludeGroups()) > 0 {
		groups = util.Intersection(groups, config.GetIncludeGroups())
	}

	if len(config.GetExcludeGroups()) > 0 {
		groups = util.Difference(groups, config.GetExcludeGroups())
	}

	for _, g := range groups {
		if lConfigs, ok := r.listeners[g]; ok {
			for _, lConfig := range lConfigs {
				listeners = append(listeners, getListenerFromConfig(ctx, lConfig))
			}
		}
	}

	return listeners
}

func (r *Registry) GetAllListeners(ctx context.Context) listener.ListenerServers {
	listeners := r.getDefaultListeners(ctx)
	for _, lConfigs := range r.listeners {
		for _, lConfig := range lConfigs {
			listeners = append(listeners, lConfig.Initializer(ctx))
		}
	}

	return listeners
}

func (r *Registry) getAllGroups() []string {
	groups := make([]string, 0)
	for group := range r.listeners {
		groups = append(groups, group)
	}

	return groups
}

func (r *Registry) getDefaultListeners(
	ctx context.Context,
) listener.ListenerServers {
	listeners := make(listener.ListenerServers, 0)
	for _, lConfig := range r.defaultListeners {
		listeners = append(listeners, getListenerFromConfig(ctx, lConfig))
	}

	return listeners
}

func (r *Registry) getListenerFromName(
	ctx context.Context,
	listenerName string,
) listener.ListenerServer {
	listenerConfig, ok := r.listenerConfigByName[listenerName]
	if !ok {
		panic(fmt.Sprintf("No listener found with name: %s", listenerName))
	}

	return getListenerFromConfig(ctx, listenerConfig)
}

func getListenerFromConfig(
	ctx context.Context,
	config listener.Config,
) listener.ListenerServer {
	l := config.Initializer(ctx)
	validateListener(l, config)

	return l
}

func validateListener(
	listener listener.ListenerServer,
	config listener.Config,
) {
	if listener.GetName() != config.Name {
		panic("listener name must be same")
	}
}
