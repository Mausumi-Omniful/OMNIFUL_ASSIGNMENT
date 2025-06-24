package newrelic

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"time"
)

type Newrelic struct {
	*newrelic.Application
}

func (nr *Newrelic) Close() error {
	nr.Shutdown(5 * time.Second)
	return nil
}
