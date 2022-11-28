package adapter

import (
	"context"
	"gochat/pkg/queue"
	"time"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

type QueueTaskAdapter struct {
	task      QueueTask          // QueueTask implementation
	queue     *queue.Queue       // message queue
	queueName string             // input queue name
	timeout   time.Duration      // QueueTask timeout
	maxWorker int                // maximum numer of goroutinue pool workers
	ctx       context.Context    // context
	cancel    context.CancelFunc // cancel function
}

func NewQueueTaskAdapter(task QueueTask, queue *queue.Queue, queueName string, timeout time.Duration, maxWorker int) *QueueTaskAdapter {
	ctx, cancel := context.WithCancel(context.Background())
	return &QueueTaskAdapter{
		task:      task,
		queue:     queue,
		queueName: queueName,
		timeout:   timeout,
		maxWorker: maxWorker,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// receive and process message constantly
func (adapter *QueueTaskAdapter) Run() {
	p, _ := ants.NewPool(adapter.maxWorker)
	defer p.Release()
	for {
		select {
		case <-adapter.ctx.Done():
			log.Info("[QueueTaskAdapter] queue task terminate")
			return
		default:
			// receive message
			message, err := adapter.queue.ReceiveMessage(adapter.queueName, 1*time.Minute)
			if err != nil {
				log.WithFields(log.Fields{
					"queueName": adapter.queueName,
				}).Errorf("[QueueTaskAdapter] failed to receive message, err = %s", err)
				continue
			}
			// process message
			p.Submit(func() {
				adapter.process(message)
			})
		}
	}
}

func (adapter *QueueTaskAdapter) Terminate() {
	adapter.cancel()
}

// process message once
func (adapter *QueueTaskAdapter) process(message queue.Message) {
	ctx, _ := context.WithTimeout(adapter.ctx, adapter.timeout)
	err := adapter.task.Run(ctx, message)
	if err != nil {
		log.WithFields(log.Fields{
			"queueName": message.QueueName,
			"data":      message.Data,
		}).Errorf("[QueueTaskAdapter] failed to process message, err = %s", err)
		return
	}
	log.WithFields(log.Fields{
		"queueName": message.QueueName,
		"data":      message.Data,
	}).Info("[QueueTaskAdapter] process message successfully")
}
