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
	maxWorker int                // maximum number of goroutine pool workers
	pool      *ants.Pool         // goroutine pool
	ctx       context.Context    // context
	cancel    context.CancelFunc // cancel function
	handler   Handler            // exception handler
}

func NewQueueTaskAdapter(task QueueTask, queue *queue.Queue, queueName string, timeout time.Duration, maxWorker int, handler Handler) *QueueTaskAdapter {
	ctx, cancel := context.WithCancel(context.Background())
	pool, _ := ants.NewPool(maxWorker)
	return &QueueTaskAdapter{
		task:      task,
		queue:     queue,
		queueName: queueName,
		timeout:   timeout,
		maxWorker: maxWorker,
		pool:      pool,
		ctx:       ctx,
		cancel:    cancel,
		handler:   handler,
	}
}

// Start receives and processes message constantly
func (adapter *QueueTaskAdapter) Start() {
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
					"taskName":  adapter.task.Name(),
				}).Errorf("[QueueTaskAdapter] failed to receive message, err = %s", err)
				continue
			}
			// process message
			adapter.pool.Submit(func() {
				adapter.process(message)
			})
		}
	}
}

func (adapter *QueueTaskAdapter) Terminate() {
	adapter.cancel()
	adapter.pool.Release()

	// print logger info
	logger, ok := adapter.handler.(*Logger)
	if ok {
		logger.Log()
	}
}

// process message once
func (adapter *QueueTaskAdapter) process(message queue.Message) {
	ctx, cancel := context.WithTimeout(adapter.ctx, adapter.timeout)
	defer cancel()

	done := make(chan struct{})

	adapter.pool.Submit(func() {
		info, status, err := adapter.task.Run(message)
		// check timeout
		adapter.checkTimeout(info, &status, err)
		// record process result
		adapter.handler.Handle(info, status, err)
		close(done)
	})

	// check timeout
	select {
	case <-done:
		return
	case <-ctx.Done():
		log.WithFields(log.Fields{
			"queueName": message.QueueName,
			"taskName":  adapter.task.Name(),
			"data":      message.Data,
		}).Errorf("[QueueTaskAdapter] process message timeout")
	}
}

func (adapter *QueueTaskAdapter) checkTimeout(info QueueTaskInfo, status *QueueTaskStatus, err error) {
	if err == nil && info.Duration > adapter.timeout {
		*status = QueueTaskStatusTimeout
	}
}
