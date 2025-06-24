package sqs

import (
	"context"
	"errors"
	"fmt"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"runtime/debug"
	"sync"
)

// Worker handles all the work
type Worker struct {
	ID          string
	taskChan    chan *[]Message
	workerPanic chan string
}

// NewWorker returns new instance of worker
func NewWorker(channel chan *[]Message, ID string, workerPanic chan string) *Worker {
	return &Worker{
		ID:          ID,
		taskChan:    channel,
		workerPanic: workerPanic,
	}
}

// Start starts the worker
func (wr *Worker) Start(ctx context.Context, wg *sync.WaitGroup, handler ISqsMessageHandler) {
	defer func() {
		if r := recover(); r != nil {
			cusErr := errors.New(fmt.Sprintf("panic occurred: %+v, stacktrace: %+v", r, string(debug.Stack())))
			newrelic.RecordEvent(newrelic.Error, map[string]interface{}{
				"Err": cusErr.Error()})
			wr.workerPanic <- wr.ID
		}
		wg.Done()
	}()

	for task := range wr.taskChan {
		err := handler.Process(ctx, task)
		if err != nil {
			log.Errorf("Sqs worker process message error: %s", err.Error())
		}
	}
}
