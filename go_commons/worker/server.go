package worker

import (
	"context"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	"github.com/omniful/go_commons/worker/configs"
	"github.com/omniful/go_commons/worker/listener"
	"github.com/omniful/go_commons/worker/registry"
)

type Server struct {
	registry  *registry.Registry
	Listeners listener.ListenerServers
}

func NewServer(listeners []listener.ListenerServer) *Server {
	return &Server{
		registry:  registry.NewRegistry(),
		Listeners: listeners,
	}
}

func NewServerFromRegistry(registry *registry.Registry) *Server {
	return &Server{
		registry:  registry,
		Listeners: make(listener.ListenerServers, 0),
	}
}

func (s *Server) RunFromConfig(
	ctx context.Context,
	config configs.ServerConfig,
) {
	s.Listeners = append(s.Listeners, s.registry.GetListenersFromConfig(ctx, config)...)
	validateListeners(s.Listeners)

	log.Infof("Starting Listeners: %s", s.Listeners.GetNames())

	s.runWithShutdownHandling(ctx)
}

func (s *Server) Run(ctx context.Context) {
	s.Listeners = append(s.Listeners, s.registry.GetAllListeners(ctx)...)
	validateListeners(s.Listeners)

	log.Infof("Starting Listeners: %s", s.Listeners.GetNames())

	s.run(ctx)
}

func (s *Server) Close() error {
	for _, server := range s.Listeners {
		server.Stop()
	}

	return nil
}

func (s *Server) runWithShutdownHandling(ctx context.Context) {
	s.run(ctx)

	shutdown.RegisterShutdownCallback("worker shutdown callback", s)
	<-shutdown.GetWaitChannel()
}

func (s *Server) run(ctx context.Context) {
	for _, server := range s.Listeners {
		go server.Start(ctx)
	}
}

func validateListeners(listeners listener.ListenerServers) {
	if len(listeners) == 0 {
		panic("no listener found to run")
	}
}
