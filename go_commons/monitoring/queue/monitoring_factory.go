package queue

import (
	"errors"
	"fmt"
)

type queueMonitoringFactory struct {
	monitors map[Monitor]Monitoring
}

type Monitor string

var _queueMonitoringFactory queueMonitoringFactory

func init() {
	_queueMonitoringFactory = queueMonitoringFactory{monitors: make(map[Monitor]Monitoring)}
}

func Register(monitor Monitor, config MonitoringConfig) error {
	_, ok := _queueMonitoringFactory.monitors[monitor]
	if ok {
		return errors.New(fmt.Sprintf("registry for %s is already done", monitor))
	}

	_queueMonitoringFactory.monitors[monitor] = NewMonitoring(config)

	return nil
}

func GetQueueMonitor(monitor Monitor) (queueMonitor Monitoring, err error) {
	queueMonitor, ok := _queueMonitoringFactory.monitors[monitor]
	if !ok {
		err = errors.New(fmt.Sprintf("Invalid Monitor Name: %s", monitor))

		return
	}

	return
}

func (m Monitor) String() string {
	return string(m)
}
