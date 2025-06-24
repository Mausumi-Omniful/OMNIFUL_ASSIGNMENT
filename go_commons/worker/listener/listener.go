package listener

import "context"

type Initializer func(ctx context.Context) ListenerServer

type ListenerServer interface {
	Start(ctx context.Context)
	Stop()
	GetName() string
}

type ListenerServers []ListenerServer

func (ls ListenerServers) GetNames() []string {
	listenerNames := make([]string, 0)

	for _, listener := range ls {
		listenerNames = append(listenerNames, listener.GetName())
	}

	return listenerNames
}
