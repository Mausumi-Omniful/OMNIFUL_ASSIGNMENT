package listener

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

type HttpListener struct {
	server      *http.Server
	serviceName string
}

func NewHttpListener(
	server *http.Server,
	serviceName string,
) ListenerServer {
	return &HttpListener{server: server, serviceName: serviceName}
}

func (l *HttpListener) Start(ctx context.Context) {
	log.Info("Stating Http Server")

	err := l.server.StartServer(l.serviceName)
	if err != nil {
		log.Errorf(fmt.Sprintf("error while starting server, err: %s", err.Error()))
		panic(err)
	}

	log.Info("Http Server started successfully")
}

func (l *HttpListener) Stop() {
	log.Info("stopping Http Server")
	if l.server == nil {
		return
	}

	err := l.server.Close()
	if err != nil {
		log.Errorf(fmt.Sprintf("error while closing server, err: %s", err.Error()))
		panic(err)
	}

	log.Info("Http Server stopped")
}

// GetName returns the name of the listener.
func (l *HttpListener) GetName() string {
	return l.serviceName
}
