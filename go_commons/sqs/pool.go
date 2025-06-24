package sqs

import (
	"context"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/util"
	"sync"
)

type Pool struct {
	name        string
	concurrency int
	collector   chan *[]Message
	wg          sync.WaitGroup
	handler     ISqsMessageHandler
	panicWorker chan string
	workers     map[string]*Worker
}

func NewPool(concurrency int, handler ISqsMessageHandler, name string) *Pool {
	return &Pool{
		name:        name,
		concurrency: concurrency,
		handler:     handler,
		panicWorker: make(chan string),
		workers:     make(map[string]*Worker, 0),
		collector:   make(chan *[]Message, util.Max(2*concurrency, 20)),
	}
}

// AddTask adds a task to the pool
func (p *Pool) AddTask(task *[]Message) {
	p.collector <- task
}

// Run runs the pool
func (p *Pool) Run(ctx context.Context) {
	for i := 1; i <= p.concurrency; i++ {
		id := uuid.New().String()
		worker := NewWorker(p.collector, id, p.panicWorker)
		p.workers[id] = worker
		p.wg.Add(1)
		go worker.Start(ctx, &p.wg, p.handler)
	}

	go p.recoveryWorker(ctx)
}

// Close closes collector channel
func (p *Pool) Close() {
	close(p.collector)
	close(p.panicWorker)
	p.wg.Wait()
}

func (p *Pool) recoveryWorker(ctx context.Context) {
	for panicWorker := range p.panicWorker {
		id := uuid.New().String()
		delete(p.workers, panicWorker)
		worker := NewWorker(p.collector, id, p.panicWorker)
		p.wg.Add(1)
		p.workers[id] = worker
		go worker.Start(ctx, &p.wg, p.handler)
		p.Monitor()
	}
}

func (p *Pool) Monitor() {
	log.Infof("Worker %v panicked........!!!! Recovering... Current worker count is %d", p.name, len(p.workers))
}
